package database

import (
	"context"
	"fmt"
	"library-api/internal/apperrors"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ApplyMigrations(pool *pgxpool.Pool, isUp bool) error {
	dirName := "../../../../migrations"
	files, err := os.ReadDir(dirName)

	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrMigrationApply, err)
	}

	filesToRead := make([]string, 0, len(files))

	for _, f := range files {
		if isUp {
			if !f.IsDir() && strings.Contains(f.Name(), ".up.sql") {
				filesToRead = append(filesToRead, f.Name())
			}
		} else {
			if !f.IsDir() && strings.Contains(f.Name(), ".down.sql") {
				filesToRead = append(filesToRead, f.Name())
			}
		}
	}

	for i := range filesToRead {
		file, err := os.ReadFile(dirName + "/" + filesToRead[i])

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
