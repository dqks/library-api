package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
)

func TestGetBook(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		cancel         bool
		wantBook       model.Book
		wantErr        error
		repositoryErr  error
		repositoryBook model.Book
		getInRepo      bool
	}{
		{
			name:   "success",
			id:     1,
			cancel: false,
			wantBook: model.Book{
				ID:        1,
				Title:     "Mock Title",
				Author:    "Mock Author",
				Year:      2000,
				Available: false,
			},
			wantErr: nil,
			repositoryBook: model.Book{
				ID:        1,
				Title:     "Mock Title",
				Author:    "Mock Author",
				Year:      2000,
				Available: false,
			},
			getInRepo: true,
		},
		{
			name:     "canceled",
			id:       1,
			wantBook: model.Book{},
			wantErr:  context.Canceled,
			cancel:   true,
		},
		{
			name:           "not found",
			id:             1,
			cancel:         false,
			wantBook:       model.Book{},
			wantErr:        apperrors.ErrNotFound,
			repositoryBook: model.Book{},
			repositoryErr:  apperrors.ErrNotFound,
			getInRepo:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.cancel {
				cancel()
			}

			repo := &MockRepo{
				err:  tt.repositoryErr,
				book: tt.repositoryBook,
			}

			s := service.Create(repo)

			book, err := s.GetBookByID(ctx, tt.id)

			if tt.getInRepo && repo.gotInside {
				if repo.id != tt.id {
					t.Fatalf("expected to call repository with id %d but got %d", tt.id, repo.id)
				}
			} else if tt.getInRepo && !repo.gotInside {
				t.Fatalf("expected not to get in repository but got")

			} else if !tt.getInRepo && repo.gotInside {
				t.Fatalf("expected to get in repository but didn't")
			}

			if err != tt.wantErr {
				t.Fatalf(
					"expected err %v but got %v",
					tt.wantErr,
					err,
				)
			}

			if tt.wantBook != book {
				t.Fatalf(
					"expected book %v but got %v",
					tt.wantBook,
					book,
				)
			}
		})
	}
}
