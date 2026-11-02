package database

import (
	"cmp"
	"context"
	"fmt"
	"library-api/internal/apperrors"
	"os"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	// dirName := "../../../../migrations"
	files, err := os.ReadDir(migrationsPath)

	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
	}

	filesToRead := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			filesToRead = append(filesToRead, f.Name())
		}
	}

	for i := range filesToRead {
		file, err := os.ReadFile(migrationsPath + "/" + filesToRead[i])

		if err != nil {
			return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
		}

		sql := string(file)

		ctx := context.Background()
		_, err = pool.Exec(ctx, sql)

		if err != nil {
			return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
		}
	}

	return nil
}

func DownMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	// dirName := "../../../../migrations"
	files, err := os.ReadDir(migrationsPath)

	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
	}

	filesToRead := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".down.sql") {
			filesToRead = append(filesToRead, f.Name())
		}
	}

	slices.SortFunc(filesToRead, func(a, b string) int {
		return cmp.Compare(b, a)
	})

	for i := range filesToRead {
		file, err := os.ReadFile(migrationsPath + "/" + filesToRead[i])

		if err != nil {
			return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
		}

		sql := string(file)

		ctx := context.Background()
		_, err = pool.Exec(ctx, sql)

		if err != nil {
			return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
		}
	}

	return nil
}
