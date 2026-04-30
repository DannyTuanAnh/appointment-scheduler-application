package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/config"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB     *sqlc.Queries
	DBPool *pgxpool.Pool
)

func InitDB() error {
	connStr := config.NewConfigDB().DB_DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Println("Error parsing DB config:", err)
		return fmt.Errorf("error parsing DB config: %w", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		log.Println("Error creating DB pool:", err)
		return fmt.Errorf("error creating DB pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close() // close pool if ping fails
		return fmt.Errorf("error pinging DB: %w", err)
	}

	// set pool and queries
	DBPool = pool
	DB = sqlc.New(DBPool)

	log.Println("Connecting to database successfully")

	return nil
}

// Close connection pool (call when app is shutting down)
func Close() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection closed")
	}
}
