package testutil

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
)

type GetBookFields struct {
	ID       int
	Ctx      context.Context
	Repo     service.Repository
	WantErr  error
	WantBook model.Book
}

func CheckGetBook(t *testing.T, f *GetBookFields) {
	book, err := f.Repo.GetBookByID(f.Ctx, f.ID)

	if !errors.Is(err, f.WantErr) {
		t.Fatalf("expected err %v but got %v", f.WantErr, err)
	}

	if book != f.WantBook {
		t.Fatalf("expected book %v but got %v", f.WantBook, book)
	}
}
