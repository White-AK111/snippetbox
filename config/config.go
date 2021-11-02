package config

import (
	"github.com/kkyr/fig"
	"log"
	"time"
)

// Config structure for settings of application
type Config struct {
	Server struct {
		ServerAddress             string        `fig:"serverAddress" default:"localhost"`                                    // address of server
		ServerPort                int           `fig:"serverPort" default:"4000"`                                            // port of server
		DSN                       string        `fig:"dsn" default:"postgres://postgres:P@ssw0rd@localhost:5432/snippetbox"` // connection string for DB
		PostgresMaxConns          int32         `fig:"postgresMaxConns" default:"8"`                                         // max connection in PG pool
		PostgresMinConns          int32         `fig:"postgresMinConns" default:"4"`                                         // min connection in PG pool
		PostgresHealthCheckPeriod time.Duration `fig:"postgresHealthCheckPeriod" default:"1"`                                // health check period in PG pool, mm
		PostgresMaxConnLifetime   time.Duration `fig:"postgresMaxConnLifetime" default:"24"`                                 // connection life in PG pool, hh
		PostgresMaxConnIdleTime   time.Duration `fig:"postgresMaxConnIdleTime" default:"30"`                                 // connection idle in PG pool, mm
		PostgresConnectTimeout    time.Duration `fig:"postgresConnectTimeout" default:"1"`                                   // connection timeout PG pool, ss
		AttackerDuration          time.Duration `fig:"attackerDuration" default:"10"`                                        // duration attack, ss
		AttackerGoroutinesCount   int           `fig:"attackerGoroutinesCount" default:"1000"`                               // count of goroutines attacker
	} `fig:"server"`
}

// Init function for initialize Config structure
func Init() (*Config, error) {
	var cfg = Config{}
	err := fig.Load(&cfg, fig.Dirs("../../config/", "../config/"), fig.File("config.yaml"))
	if err != nil {
		log.Fatalf("can't load configuration file: %s", err)
		return nil, err
	}

	return &cfg, err
}