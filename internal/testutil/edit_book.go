package testutil

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"slices"
	"testing"
)

type EditFields struct {
	Err            error
	Book           model.Book
	CtxGet         context.Context
	Repo           service.BookRepository
	Cancel         bool
	BooksBefore    []model.Book
	WantRepoChange bool
	WantErr        error
	WantBook       model.Book
}

func CheckEdit(t *testing.T, f *EditFields) {
	if !errors.Is(f.Err, f.WantErr) {
		t.Fatalf("expected err %v but got %v", f.WantErr, f.Err)
	}

	if f.Book != f.WantBook {
		t.Fatalf("expected book %v but got %v", f.WantBook, f.Book)
	}

	booksAfter, err := f.Repo.GetBooks(f.CtxGet, repository.GetBooksPayload{})

	if err != nil {
		t.Fatal("failed to get books")
	}

	if len(booksAfter) != len(f.BooksBefore) {
		t.Fatal("got repo unexpectedly changed")
	}

	for i := range booksAfter {
		if booksAfter[i] != f.BooksBefore[i] && !f.WantRepoChange {
			t.Fatalf(
				"got book %v changed unexpectedly %v",
				f.BooksBefore[i],
				booksAfter[i],
			)
		} else if f.WantRepoChange &&
			booksAfter[i] != f.BooksBefore[i] &&
			((f.WantBook.ID != booksAfter[i].ID) || (f.WantBook.ID != f.BooksBefore[i].ID)) {
			t.Fatalf(
				"got book %v changed unexpectedly %v",
				f.BooksBefore[i],
				booksAfter[i],
			)
		}
	}

	if !f.Cancel && f.WantErr == nil {
		index := slices.Contains(booksAfter, f.Book)

		if !index {
			t.Fatalf("failed to find book in repository %v", f.Book)
		}
	}
}
