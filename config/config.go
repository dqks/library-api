package config

import (
	"fmt"
	"library-api/internal/apperrors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN     string
	TestDSN string
	Port    int
}

func Get(prod bool) (*Config, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	env := findDotEnv(currentDir)
	dsn, err := GetEnvVar(env, "DATABASE_URL")

	if err != nil {
		return nil, err
	}

	var testDsn string
	if !prod {
		testDsn, err = GetEnvVar(env, "DATABASE_TEST_URL")

		if err != nil {
			return nil, err
		}
	}

	return &Config{
		DSN:     dsn,
		TestDSN: testDsn,
		Port:    8080,
	}, nil
}

func GetEnvVar(envPath, envVar string) (string, error) {
	dsnGlobalEnv, exists := os.LookupEnv(envVar)
	fmt.Println(dsnGlobalEnv, "dsnGlobalEnv 50 line")
	fmt.Println(exists, "exists 51 line")
	if exists && dsnGlobalEnv == "" {
		return "", apperrors.ErrBadEnvironmentConfiguration
	}

	if err := godotenv.Load(envPath); err != nil {
		// Если путь неправильный .env и глобальный есть
		// то возвращаем его
		fmt.Println(dsnGlobalEnv, "dsnGlobalEnv 59 line")
		if dsnGlobalEnv != "" {
			return dsnGlobalEnv, nil
		}

		// Если путь неправильный .env и глобального нет
		// то возвращаем ошибку
		return "", fmt.Errorf("%w: %v", apperrors.ErrInvalidFileLoad, err)
	}

	dsnInnerEnv := os.Getenv(envVar)

	// Если обе переменные загружены возвращаем глобальную
	// Если глобальная загружена а .env не загружена возвращаем глобальную
	// иначе возвращаем .env
	fmt.Println(dsnGlobalEnv, "dsnGlobalEnv 74 line")
	fmt.Println(dsnInnerEnv, "dsnGlobalEnv 75 line")

	if dsnGlobalEnv != "" {
		return dsnGlobalEnv, nil
	} else if dsnInnerEnv != "" {
		return dsnInnerEnv, nil
	}

	return "", apperrors.ErrBadEnvVariableName
}
func findDotEnv(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir + "/.env"
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
