package testutil

import (
	"context"
	"library-api/config"
	"library-api/internal/database"
	"library-api/internal/repository/postgres"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePostgresTestRepo(t *testing.T, dsnPath string) (*postgres.BookRepository, *pgxpool.Pool) {
	config, err := config.Get(false)

	if err != nil {
		t.Fatalf("failed to get dsn %v", err)
	}

	ctx := context.Background()

	pool, err := database.CreatePool(ctx, config.DSN)

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

func SeedBooks(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()

	_, err := ClearBooks(ctx, pool)

	if err != nil {
		t.Fatalf("failed to clear books %v", err)
	}

	_, err = pool.Exec(
		ctx,
		`
		INSERT INTO books (id, title, author, year, available)
		VALUES 
			(1, 'Title 1', 'Author 1', 2000, true),
			(2, 'Title 2', 'Author 2', 2000, true),
			(3, 'Title 3', 'Author 3', 2000, false),
			(4, 'Title 4', 'Author 4', 2000, true),
			(5, 'Title 5', 'Author 5', 2000, true);
		`,
	)

	if err != nil {
		t.Fatalf("failed to seed into books %v", err)
	}

	_, err = pool.Exec(
		ctx,
		`
		SELECT setval('books_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM books), false);
		`,
	)

	if err != nil {
		t.Fatalf("failed to set id counter %v", err)
	}
}

func ClearUsers(ctx context.Context, pool *pgxpool.Pool) (pgconn.CommandTag, error) {
	return pool.Exec(
		ctx,
		"DELETE FROM users;",
	)
}

func SeedUsers(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()

	_, err := ClearUsers(ctx, pool)

	if err != nil {
		t.Fatalf("failed to clear users %v", err)
	}

	_, err = pool.Exec(
		ctx,
		`
		INSERT INTO users (id, name)
		VALUES 
			(1, 'user 1'),
			(2, 'user 2'),
			(3, 'user 3'),
			(4, 'user 4'),
			(5, 'user 5')
		`,
	)

	if err != nil {
		t.Fatalf("failed to seed into users %v", err)
	}

	_, err = pool.Exec(
		ctx,
		`
		SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM users), false);
		`,
	)

	if err != nil {
		t.Fatalf("failed to set id counter %v", err)
	}
}

func ClearLoans(ctx context.Context, pool *pgxpool.Pool) (pgconn.CommandTag, error) {
	return pool.Exec(
		ctx,
		"DELETE FROM loans;",
	)
}

func SeedLoans(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()

	_, err := ClearLoans(ctx, pool)

	if err != nil {
		t.Fatalf("failed to clear loans %v", err)
	}

	_, err = pool.Exec(
		ctx,
		`
		INSERT INTO loans (id, user_id, book_id, loan_date, return_date)
		VALUES 
			(1, 1, 1, '2026-10-29', null),
			(2, 1, 2, '2026-10-29', '2026-11-02')
		`,
	)

	if err != nil {
		t.Fatalf("failed to seed into loans %v", err)
	}

	_, err = pool.Exec(
		ctx,
		`
		SELECT setval('loans_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM loans), false);
		`,
	)

	if err != nil {
		t.Fatalf("failed to set id counter %v", err)
	}
}
