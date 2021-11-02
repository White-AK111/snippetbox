package main

import (
	"context"
	"fmt"
	"github.com/White-AK111/snippetbox/config"
	"github.com/White-AK111/snippetbox/pkg/models/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// app struct for web-application
type appAttack struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

// AttackResults struct for get result of attack
type AttackResults struct {
	Duration         time.Duration
	GoroutinesCount  int
	QueriesPerformed uint64
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Can't load configuration file: %s", err)
	}

	// init logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// init app
	app := &appAttack{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	if err = app.initPgServer(cfg); err != nil {
		log.Fatalf("Can't connect to DB: %s", err)
	}
	defer app.snippets.DB.Close()

	//lastSnippets, err := app.snippets.LatestSnippets(1, 10)
	//if err != nil {
	//log.Fatalf("Error on func Latest: %s", err)
	//}

	//for _, snippet := range lastSnippets {
	//	fmt.Println(*snippet)
	//}

	//testQ, err := app.snippets.GetUserByLogin("mike111")
	//fmt.Println(testQ)

	attackResults := app.attack(cfg)
	fmt.Printf("Result of attack: %+v\n", attackResults)
}

// initPgServer method init connections to postgres
func (a *appAttack) initPgServer(cfg *config.Config) error {
	ctx := context.Background()

	cfgPg, err := pgxpool.ParseConfig(cfg.Server.DSN)
	if err != nil {
		log.Fatal(err)
	}

	cfgPg.MaxConns = cfg.Server.PostgresMaxConns
	cfgPg.MinConns = cfg.Server.PostgresMinConns
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

// attack method for attack postgres concurrently
func (a *appAttack) attack(cfg *config.Config) AttackResults {
	var queries uint64
	var wg sync.WaitGroup

	attacker := func(stopAt time.Time) {
		for {
			_, err := a.snippets.GetUserByLogin("mike111")
			if err != nil {
				log.Printf("Error on GetUserByLogin attack: %s", err)
			}

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	startAt := time.Now()
	stopAt := startAt.Add(cfg.Server.AttackerDuration * time.Second)

	for i := 0; i < cfg.Server.AttackerGoroutinesCount; i++ {
		wg.Add(1)
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		GoroutinesCount:  cfg.Server.AttackerGoroutinesCount,
		QueriesPerformed: queries,
	}
}
