package testutil

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"testing"
)

type DeleteBookFields struct {
	ID          int
	Ctx         context.Context
	Cancel      bool
	BooksBefore []model.Book
	Repo        service.BookRepository
	WantErr     error
}

func CheckDeleteBook(t *testing.T, f *DeleteBookFields) {
	err := f.Repo.DeleteBookByID(f.Ctx, f.ID)

	if !errors.Is(err, f.WantErr) {
		t.Fatalf("expected err %v but got %v", f.WantErr, err)
	}
	ctxGet := context.WithoutCancel(context.Background())
	booksAfter, err := f.Repo.GetBooks(ctxGet, repository.GetBooksPayload{})
	if err != nil {
		t.Fatal("failed to get books")
	}

	if f.WantErr != nil {
		if len(booksAfter) != len(f.BooksBefore) {
			t.Fatalf(
				"got len of books %d but before delete it was %d",
				len(booksAfter),
				len(f.BooksBefore),
			)
		}

		for i := 0; i < len(booksAfter); i++ {
			if booksAfter[i] != f.BooksBefore[i] {
				t.Fatalf("got changes in books after delete")
			}
		}
	}

	if f.WantErr == nil && !f.Cancel {
		if len(booksAfter) != len(f.BooksBefore)-1 {
			t.Fatalf(
				"got len of books %d but before delete it was %d",
				len(booksAfter),
				len(f.BooksBefore),
			)
		}

		for _, b := range booksAfter {
			if b.ID == f.ID {
				t.Fatalf("expected to delete book %v but got it in books", b)
			}
		}
	}
}
