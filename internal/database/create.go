package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func CreatePool(envPath, envVar string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	dsn := os.Getenv(envVar)
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		pool.Close()
		return nil, err

	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err

	}

	return pool, nil
}
