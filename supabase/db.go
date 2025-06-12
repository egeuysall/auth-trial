package supabase

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func Connect() *pgxpool.Pool {
	connStr := os.Getenv("POSTGRES_URL")

	if connStr == "" {
		log.Fatal("POSTGRES_URL is not set")
	}

	config, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		log.Fatalf("Failed to parse db url: %s", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 30 * time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Fatalf("Failed to create database pool: %s", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %s", err)
	}

	log.Println("Successfully connected Supabase database")
	return pool
}
