package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	/*Database details*/
	databaseURL := "postgres://user:password@localhost:5432/postgres?sslmode=disable"

	ctx := context.Background()

	/*Connect to database*/
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}

	Pool = pool

	log.Println("Connected to PostgreSQL")
	/*Prepare queries*/
	queries := []string{
		`CREATE TABLE IF NOT EXISTS system_logs (
			host_id TEXT NOT NULL,
			recorded_at TIMESTAMPTZ NOT NULL,
			cpu_usage DOUBLE PRECISION,
			total_ram BIGINT,
			free_ram  BIGINT,
			used_ram  DOUBLE PRECISION,
			available_ram BIGINT,
			total_disk BIGINT,
			used_disk  BIGINT,
			free_disk  BIGINT,
			PRIMARY KEY (host_id, recorded_at)
		);`,

		`CREATE TABLE IF NOT EXISTS service_logs (
			host_id TEXT NOT NULL,
			service_name TEXT NOT NULL,
			recorded_at TIMESTAMPTZ NOT NULL,
			status TEXT,
			cpu_usage DOUBLE PRECISION,
			ram_usage  DOUBLE PRECISION,
			PRIMARY KEY (host_id, recorded_at, service_name)
		);`,
	}
	/*Create tables*/
	for _, q := range queries {
		if _, err := Pool.Exec(ctx, q); err != nil {
			log.Fatalf("Failed to create table: %v\n", err)
		}
	}

	log.Println("Database schema initialized successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
