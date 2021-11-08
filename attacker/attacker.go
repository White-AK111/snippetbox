package main

import (
	"context"
	"github.com/White-AK111/snippetbox/config"
	"github.com/White-AK111/snippetbox/pkg/models/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// workerPool struct for management goroutines
type workerPool struct {
	wg            sync.WaitGroup
	resultChan    chan uint64
	semaphoreChan chan struct{}
	ctx           context.Context
}

// newWorkerPool method initialize new WorkerPool, return *workerPool
func newWorkerPool(N uint) *workerPool {
	return &workerPool{
		wg:            sync.WaitGroup{},
		resultChan:    make(chan uint64, N),
		semaphoreChan: make(chan struct{}, N),
		ctx:           context.Background(),
	}
}

// app struct for web-application
type appAttack struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

// AttackResults struct for get result of attack
type AttackResults struct {
	MaxConn          uint32
	MinConn          uint32
	Duration         time.Duration
	GoroutinesCount  uint
	QueriesPerformed uint64
	QPS              uint64
}

func main() {
	// init logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cfg, err := config.Init()
	if err != nil {
		errorLog.Fatalf("Can't load configuration file: %s\n", err)
	}

	// init app
	app := &appAttack{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	if err = app.initPgServer(cfg); err != nil {
		errorLog.Fatalf("Can't connect to DB: %s\n", err)
	}
	defer app.snippets.DB.Close()

	wp := newWorkerPool(cfg.Server.AttackerGoroutinesCount)

	// !!! fill data in DB by data.sql before use this function !!!
	attackResults := app.attack(cfg, wp)
	infoLog.Printf("Result of attack: %+v\n", attackResults)
}

// initPgServer method init connections to postgres
func (a *appAttack) initPgServer(cfg *config.Config) error {
	ctx := context.Background()

	cfgPg, err := pgxpool.ParseConfig(cfg.Server.DSN)
	if err != nil {
		return err
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

// attack method for attack postgres
func (a *appAttack) attack(cfg *config.Config, wp *workerPool) AttackResults {
	var queries uint64

	startAt := time.Now()
	stopAt := startAt.Add(cfg.Server.AttackerDuration * time.Second)
	t := time.NewTimer(cfg.Server.AttackerDuration * time.Second)

	for {
		select {
		case <-wp.ctx.Done():
			{
				wp.wg.Wait()
				return AttackResults{
					MaxConn:          cfg.Server.PostgresMaxConns,
					MinConn:          cfg.Server.PostgresMinConns,
					Duration:         time.Now().Sub(startAt),
					GoroutinesCount:  cfg.Server.AttackerGoroutinesCount,
					QueriesPerformed: queries,
					QPS:              uint64(math.Round(float64(queries) / float64(cfg.Server.AttackerDuration))),
				}
			}
		case wp.semaphoreChan <- struct{}{}:
			{
				wp.wg.Add(1)
				go attacker(cfg, "GetUserByLogin", stopAt, a, wp, &queries)
			}
		case wp.semaphoreChan <- struct{}{}:
			{
				wp.wg.Add(1)
				go attacker(cfg, "LatestSnippets", stopAt, a, wp, &queries)
			}
		case wp.semaphoreChan <- struct{}{}:
			{
				wp.wg.Add(1)
				go attacker(cfg, "GetNotSendedNotifications", stopAt, a, wp, &queries)
			}
		case <-t.C:
			{
				wp.wg.Wait()
				return AttackResults{
					MaxConn:          cfg.Server.PostgresMaxConns,
					MinConn:          cfg.Server.PostgresMinConns,
					Duration:         time.Now().Sub(startAt),
					GoroutinesCount:  cfg.Server.AttackerGoroutinesCount,
					QueriesPerformed: queries,
					QPS:              uint64(math.Round(float64(queries) / float64(cfg.Server.AttackerDuration))),
				}
			}
		default:
			wp.wg.Wait()
		}
	}
}

// attacker function for attack postgres concurrently
func attacker(cfg *config.Config, query string, stopAt time.Time, a *appAttack, wp *workerPool, queries *uint64) {
	defer wp.wg.Done()
	for {
		if time.Now().Before(stopAt) {
			switch query {
			case "GetUserByLogin":
				{
					prefix := cfg.Server.AttackerPrefixUserInDB
					userId := getRandomInt(int(cfg.Server.AttackerCountUserInDB))
					login := prefix + strconv.Itoa(userId)
					_, err := a.snippets.GetUserByLogin(login)
					if err != nil {
						a.errorLog.Printf("Error on GetUserByLogin attack: %s login: %s \n", err, login)
					}
					atomic.AddUint64(queries, 1)
				}
			case "LatestSnippets":
				{
					userId := getRandomInt(int(cfg.Server.AttackerCountUserInDB))
					limit := getRandomInt(int(cfg.Server.AttackerLimitSelectCount))
					_, err := a.snippets.LatestSnippets(uint(userId), uint(limit))
					if err != nil {
						log.Printf("Error on GetUserByLogin attack: %s login: %s limit: %s \n", err, userId, limit)
					}
					atomic.AddUint64(queries, 1)
				}
			case "GetNotSendedNotifications":
				{
					_, err := a.snippets.GetNotSendedNotifications()
					if err != nil {
						a.errorLog.Printf("Error on GetNotSendedNotifications attack: %s \n", err)
					}
					atomic.AddUint64(queries, 1)
				}
			}
		} else {
			return
		}
	}
}

// getRandomInt get random int, where N = {1..N)
func getRandomInt(n int) int {
	rand.Seed(time.Now().UnixNano())
	rInt := rand.Intn(n)
	if rInt == 0 {
		rInt++
	}
	return n
}
