package main

import (
	"context"
	"github.com/White-AK111/snippetbox/config"
	"github.com/White-AK111/snippetbox/pkg/models/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// app struct for web-application
type app struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Can't load configuration file: %s", err)
	}

	// init logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// init templates cache
	templateCache, err := newTemplateCache("../../ui/html/")
	if err != nil {
		log.Fatalf("Can't init templates cache: %s", err)
	}

	// init app
	app := &app{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
	}

	if err = app.initPgServer(cfg); err != nil {
		log.Fatalf("Can't connect to DB: %s", err)
	}
	defer app.snippets.DB.Close()

	connStr := cfg.Server.ServerAddress + ":" + strconv.Itoa(int(cfg.Server.ServerPort))

	// init server
	srv := &http.Server{
		Addr:           connStr,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	infoLog.Printf("Start server on %s", connStr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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

	a.snippets = &postgres.SnippetModel{DB: dbPool, CTX: ctx}

	return nil
}
