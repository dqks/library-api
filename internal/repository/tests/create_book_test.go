package repository_test

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
	"testing"
)

func createTestRepo() *repository.BookRepository {
	return repository.CreateRepo(
		3,
		[]model.Book{
			{
				ID:        1,
				Title:     "1984",
				Author:    "George Orwell",
				Year:      uint16(1956),
				Available: true,
			},
			{
				ID:        2,
				Title:     "Kallocain",
				Author:    "Karin Boye",
				Year:      uint16(1937),
				Available: false,
			},
		},
	)
}

func TestCreateBook(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		payload   repository.CreateBookPayload
		wantError error
		wantBook  model.Book
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
