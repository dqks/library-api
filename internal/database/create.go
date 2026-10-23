package database

import (
	"context"
	"library-api/internal/apperrors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func GetDsn(envPath, envVar string) (string, error) {
	if err := godotenv.Load(envPath); err != nil {
		return "", err
	}

	dsn := os.Getenv(envVar)

	if dsn == "" {
		return "", apperrors.ErrBadVariableName
	}

	return dsn, nil
}

func CreatePool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err

	}

	return pool, nil
}
