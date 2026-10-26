package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpMigrations(pool *pgxpool.Pool) error {
	dirName := "../../../../migrations"
	files, err := os.ReadDir(dirName)

	if err != nil {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	filesToRead := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), ".up.sql") {
			filesToRead = append(filesToRead, f.Name())
		}
	}

	for i := range filesToRead {
		file, err := os.ReadFile(dirName + "/" + filesToRead[i])

		if err != nil {
			return fmt.Errorf("ошибка применения миграций: %w", err)
		}

		sql := string(file)

		ctx := context.Background()
		_, err = pool.Exec(ctx, sql)

		if err != nil {
			return fmt.Errorf("ошибка применения миграций: %w", err)
		}
	}

	return nil
}

func DownMigrations(pool *pgxpool.Pool) error {
	dirName := "../../../../migrations"
	files, err := os.ReadDir(dirName)

	if err != nil {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	filesToRead := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), ".down.sql") {
			filesToRead = append(filesToRead, f.Name())
		}
	}

	for i := range filesToRead {
		file, err := os.ReadFile(dirName + "/" + filesToRead[i])

		if err != nil {
			return fmt.Errorf("ошибка применения миграций: %w", err)
		}

		sql := string(file)

		ctx := context.Background()
		_, err = pool.Exec(ctx, sql)

		if err != nil {
			return fmt.Errorf("ошибка применения миграций: %w", err)
		}
	}

	return nil
}
