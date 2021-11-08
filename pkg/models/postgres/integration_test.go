//go:build integration
// +build integration

package postgres

import (
	"context"
	"github.com/White-AK111/snippetbox/config"
	"github.com/White-AK111/snippetbox/pkg/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"log"
	"net"
	"testing"
	"time"
)

// app struct for test
type app struct {
	snippets *SnippetModel
}

var appTest = &app{}

// TestIntegrationConnectionToPG test connect to DB
func TestIntegrationConnectionToPG(t *testing.T) {
	// init config
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Can't load configuration file: %s\n", err)
	}

	if err = appTest.initPgServer(cfg); err != nil {
		log.Fatalf("Can't connect to DB: %s", err)
	}
}

// initPgServer method init connections to postgres
func (a *app) initPgServer(cfg *config.Config) error {
	ctx := context.Background()

	cfgPg, err := pgxpool.ParseConfig(cfg.Server.DSN)
	if err != nil {
		log.Fatal(err)
	}

	cfgPg.MaxConns = int32(cfg.Server.PostgresMaxConns)
	cfgPg.MinConns = int32(cfg.Server.PostgresMinConns)
	cfgPg.HealthCheckPeriod = cfg.Server.PostgresHealthCheckPeriod * time.Minute
	cfgPg.MaxConnLifetime = cfg.Server.PostgresMaxConnLifetime * time.Hour
	cfgPg.MaxConnIdleTime = cfg.Server.PostgresMaxConnIdleTime * time.Minute
	cfgPg.ConnConfig.ConnectTimeout = cfg.Server.PostgresConnectTimeout * time.Second

	cfgPg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfgPg.HealthCheckPeriod,
		Timeout:   cfgPg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbPool, err := pgxpool.ConnectConfig(ctx, cfgPg)
	if err != nil {
		return err
	}

	a.snippets = &SnippetModel{DB: dbPool, CTX: ctx}

	return nil
}

// TestIntegrationGetUserByLogin test queries InsertUser and GetUserByLogin
func TestIntegrationGetUserByLogin(t *testing.T) {
	defer appTest.snippets.DB.Close()

	tests := []struct {
		name    string
		app     *app
		ctx     context.Context
		login   string
		prepare func(*testing.T)
		check   func(*testing.T, *models.User, error)
	}{
		{
			name:  "success",
			app:   appTest,
			ctx:   context.Background(),
			login: "TestUser100500",
			prepare: func(t *testing.T) {
				user := models.User{
					Name:           "TestUser100500",
					Login:          "TestUser100500",
					Email:          "TestUser100500@email.com",
					HashedPassword: pgtype.Bytea{Bytes: []byte("P@ssw0rd"), Status: pgtype.Present},
					Created:        pgtype.Timestamptz{Time: time.Now(), Status: pgtype.Present},
					Confirmed:      true,
				}
				_, err := appTest.snippets.InsertUser(&user)
				require.NoError(t, err)
			},
			check: func(t *testing.T, user *models.User, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, *user)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(t)
			user, err := tt.app.snippets.GetUserByLogin(tt.login)
			tt.check(t, user, err)
		})
	}
}
