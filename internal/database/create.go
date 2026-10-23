package database

import (
	"context"
	"library-api/internal/apperrors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func GetDsn(envPath, envVar string) (string, error) {
	dsnGlobalEnv, exists := os.LookupEnv(envVar)
	if exists && dsnGlobalEnv == "" {
		return "", apperrors.ErrBadEnvironmentConfiguration
	}

	if err := godotenv.Load(envPath); err != nil {
		// Если путь неправильный .env то возвращаем и глобальный есть
		// то возвращаем его
		if dsnGlobalEnv != "" {
			return dsnGlobalEnv, nil
		}

		// Если путь неправильный .env и глобального нет
		// то возвращаем ошибку
		return "", apperrors.ErrInvalidFileLoad
	}

	dsnInnerEnv := os.Getenv(envVar)

	// Если обе переменные загружены возвращаем глобальную
	// Если глобальная загружена а .env не загружена возвращаем глобальную
	// иначе возвращаем .env
	if dsnGlobalEnv != "" {
		return dsnGlobalEnv, nil
	} else if dsnInnerEnv != "" {
		return dsnInnerEnv, nil
	}

	return "", apperrors.ErrBadEnvVariableName
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
