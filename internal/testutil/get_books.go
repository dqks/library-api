package testutil

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"testing"
)

type GetBooksFields struct {
	Err       error
	Books     []model.Book
	Cancel    bool
	Payload   repository.GetBooksPayload
	Ctx       context.Context
	Repo      service.Repository
	WantErr   error
	WantBooks []model.Book
}

func CheckGetBooks(t *testing.T, f *GetBooksFields) {
	if !errors.Is(f.Err, f.WantErr) {
		t.Fatalf("expected err %v but got %v", f.WantErr, f.Err)
	}

	if len(f.Books) != len(f.WantBooks) {
		t.Fatalf("expected len of books %d but got %d", len(f.WantBooks), len(f.Books))
	}

	for i := range f.WantBooks {
		if f.Books[i] != f.WantBooks[i] {
			t.Fatalf("expected books %v \n but got %v", f.WantBooks, f.Books)
		}
	}

	for i := range f.WantBooks {
		f.Books[i].ID = f.Books[i].ID + 1
	}

	if !f.Cancel {
		booksAfter, err := f.Repo.GetBooks(f.Ctx, f.Payload)

		if err != nil {
			t.Fatalf("failed to get books")
		}

		for i := range f.WantBooks {
			if f.Books[i] == booksAfter[i] {
				t.Fatal("repo chnaged expectedly")
			}
		}
	}
}
