package database

import (
	"context"
	"library-api/internal/apperrors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func GetDsn(envPath, envVar string) (string, error) {
	dsnGlobalEnv := os.Getenv(envVar)

	if err := godotenv.Load(envPath); err != nil {
		// Если путь неправильный .env то возвращаем и глобальный есть
		// то возвращаем его
		if dsnGlobalEnv != "" {
			return "", apperrors.ErrBadEnvVariable
		}

		// Если путь неправильный .env то возвращаем и глобальный есть
		// то возвращаем его
		return "", err
	}

	dsnInnerEnv := os.Getenv(envVar)

	// Если обе переменные загружены возвращаем глобальную
	// Если глобальная загружена а .env не загружена возвращаем глобальную
	// иначе возвращаем .env
	if (dsnGlobalEnv != "" && dsnInnerEnv != "") || (dsnGlobalEnv != "" && dsnInnerEnv == "") {
		return dsnGlobalEnv, nil
	} else {
		return dsnInnerEnv, nil
	}
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
