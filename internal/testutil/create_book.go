package testutil

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"testing"
)

type CreateFields struct {
	Title        string
	Author       string
	Year         uint16
	Available    bool
	Ctx          context.Context
	Repo         service.Repository
	WantError    error
	WantBook     model.Book
	WantToAppear bool
	WantLen      int
}

func CheckCreateBook(t *testing.T, f *CreateFields) {
	book, err := f.Repo.CreateBook(f.Ctx, repository.CreateBookPayload{
		Title:     &f.Title,
		Author:    &f.Author,
		Year:      &f.Year,
		Available: &f.Available,
	})

	if !errors.Is(err, f.WantError) {
		t.Fatalf("got error %v but expected %v", err, f.WantError)
	}

	if book != f.WantBook {
		t.Fatalf("got %v but expected %v", book, f.WantBook)
	}

	ctxGet := context.WithoutCancel(context.Background())

	books, err := f.Repo.GetBooks(ctxGet, repository.GetBooksPayload{})

	if err != nil {
		t.Fatalf("failed to get books %v", err)
	}

	if len(books) != f.WantLen {
		t.Fatalf("got len %d but expected %d", len(books), f.WantLen)
	}

	foundBook := false
	for i := range books {
		if books[i] == book {
			foundBook = true
		}
	}
	if (!foundBook && f.WantToAppear) || (foundBook && !f.WantToAppear) {
		t.Fatalf("new book %v didn't appear on book list", book)
	}
}
