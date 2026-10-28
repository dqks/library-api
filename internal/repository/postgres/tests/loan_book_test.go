package postgres_tests

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/database"
	"testing"
)

type Loan struct {
	bookId,
	userID int
}

func TestLoanBook(t *testing.T) {
	// ID 1000 == not found
	tests := []struct {
		name    string
		bookID  int
		userID  int
		cancel  bool
		wantErr error
	}{
		{
			name:    "success",
			bookID:  1,
			userID:  1,
			wantErr: nil,
		},
		{
			name:    "book not found",
			bookID:  1000,
			userID:  1,
			wantErr: apperrors.ErrBookUnavailable,
		},
		{
			name:    "user not found",
			bookID:  1,
			userID:  1000,
			wantErr: apperrors.ErrInternal,
		},
		{
			name:    "user not found, book found",
			bookID:  2,
			userID:  1000,
			wantErr: apperrors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, pool := CreatePostgresTestRepo(t)

			t.Cleanup(func() {
				pool.Close()
			})

			err := database.UpMigrations(pool)
			if err != nil {
				t.Fatalf("migration's apply error: %v", err)
			}

			t.Cleanup(func() {
				err := database.DownMigrations(pool)
				if err != nil {
					t.Errorf("migration's apply error: %v", err)
				}
			})

			SeedBooks(t, pool)
			SeedUsers(t, pool)

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(func() {
				cancel()
			})

			if tt.cancel {
				cancel()
			}

			id, err := repo.LoanBook(ctx, tt.bookID, tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected err to be %v but got %v",
					tt.wantErr,
					err,
				)
			}

			// Зафиксировали, что если пользователь не найден,
			// то все равно у книги available будет false
			if tt.wantErr == nil {
				ctxQuery := context.Background()
				book, err := repo.GetBookByID(ctxQuery, tt.bookID)

				if err != nil {
					t.Fatalf(
						"failed to find book with id %d",
						tt.bookID,
					)
				}

				if tt.userID == 1000 {
					if book.Available != true {
						t.Fatalf(
							"expected book's available field to be true but got %v",
							book.Available,
						)
					}
				} else {
					if book.Available != false {
						t.Fatalf(
							"expected book's available field to be false but got %v",
							book.Available,
						)
					}
				}
			}

			if tt.wantErr == nil {
				ctxQuery := context.Background()

				var loan Loan
				err = pool.QueryRow(
					ctxQuery,
					`
				SELECT book_id, user_id FROM loans WHERE id = $1
				`,
					id,
				).Scan(
					&loan.bookId,
					&loan.userID,
				)

				if loan.bookId != tt.bookID {
					t.Fatalf(
						"expected loan's book id to be %d but got %d",
						tt.bookID,
						loan.bookId,
					)
				}

				if loan.userID != tt.userID {
					t.Fatalf(
						"expected loan's user id to be %d but got %d",
						tt.userID,
						loan.userID,
					)
				}
			}
		})
	}
}
