package pgrepo

import "github.com/jackc/pgx/v5/pgxpool"

type BookRepository struct {
	pool *pgxpool.Pool
}

func Create(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		pool: pool,
	}
}
