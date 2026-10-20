package postgres_tests

import (
	"context"
	"library-api/internal/database"
	"library-api/internal/repository/postgres"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePostgresTestRepo(t *testing.T) (*postgres.BookRepository, *pgxpool.Pool) {
	pool, err := database.CreatePool("../../../../.env", "DATABASE_TEST_URL")

	if err != nil {
		t.Fatalf("failed to create repository %v", err)
	}

	return postgres.Create(pool), pool
}

func ClearBooks(ctx context.Context, pool *pgxpool.Pool) (pgconn.CommandTag, error) {
	return pool.Exec(
		ctx,
		"DELETE FROM books;",
	)
}

func SeedBooks(t *testing.T, pool *pgxpool.Pool) (pgconn.CommandTag, error) {
	ctx := context.Background()

	_, err := ClearBooks(ctx, pool)

	if err != nil {
		t.Fatalf("failed to clear books %v", err)
	}

	tag, err := pool.Exec(
		ctx,
		`
		INSERT INTO books (id, title, author, year, available)
		VALUES 
			(1, 'Title 1', 'Author 1', 2000, true),
			(2, 'Title 2', 'Author 2', 2000, true),
			(3, 'Title 3', 'Author 3', 2000, true),
			(4, 'Title 4', 'Author 4', 2000, true),
			(5, 'Title 5', 'Author 5', 2000, true);
		`,
	)

	if err != nil {
		t.Fatalf("failed to seed into books %v", err)
	}

	tag, err = pool.Exec(
		ctx,
		`
		SELECT setval('books_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM books), false);
		`,
	)

	if err != nil {
		t.Fatalf("failed to set id counter %v", err)
	}

	return tag, err
}
