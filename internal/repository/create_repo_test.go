package repository

import (
	"library-api/internal/model"
	"testing"
)

func TestCreateRepo(t *testing.T) {
	tests := []struct {
		name        string
		books       []model.Book
		nextID      int
		wantBooks   []model.Book
		wantRepoNil bool
		wantNextID  int
	}{
		{
			name:   "success",
			nextID: 3,
			books: []model.Book{
				{
					ID:        1,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        2,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        2,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantRepoNil: false,
			wantNextID:  3,
		},
		{
			name:   "nil repo, ID id more than nextID",
			nextID: 3,
			books: []model.Book{
				{
					ID:        6,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        2,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantBooks:   []model.Book{},
			wantRepoNil: true,
			wantNextID:  0,
		},
		{
			name:   "nil repo, nextID is 0",
			nextID: 0,
			books: []model.Book{
				{
					ID:        6,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        2,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantBooks:   []model.Book{},
			wantRepoNil: true,
			wantNextID:  0,
		},
		{
			name:   "nil repo, nextID is -1",
			nextID: -1,
			books: []model.Book{
				{
					ID:        6,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        2,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantBooks:   []model.Book{},
			wantRepoNil: true,
			wantNextID:  0,
		},
		{
			name:   "nil repo, same id",
			nextID: 3,
			books: []model.Book{
				{
					ID:        1,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        1,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantBooks:   []model.Book{},
			wantRepoNil: true,
			wantNextID:  0,
		},
		{
			name:   "nil repo, id is -2",
			nextID: 3,
			books: []model.Book{
				{
					ID:        1,
					Title:     "Test Title 1",
					Author:    "Test Author 1",
					Year:      2003,
					Available: true,
				},
				{
					ID:        -2,
					Title:     "Test Title 2",
					Author:    "Test author 2",
					Year:      2005,
					Available: false,
				},
			},
			wantBooks:   []model.Book{},
			wantRepoNil: true,
			wantNextID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := CreateRepo(tt.nextID, tt.books)

			if repo != nil && tt.wantRepoNil || repo == nil && !tt.wantRepoNil {
				if tt.wantRepoNil {
					t.Fatalf("got repo %v but expected nil", repo)
				} else {
					t.Fatalf("got repo %v but expected not nil", repo)
				}
			}

			if !tt.wantRepoNil {

				if len(repo.books) != len(tt.wantBooks) {
					t.Fatalf(
						"expected len of repo's books %d but got %d",
						len(tt.wantBooks),
						len(repo.books),
					)
				}

				for i := range tt.wantBooks {
					if tt.wantBooks[i] != repo.books[i] {
						t.Fatalf(
							"expected book %v but got %v",
							tt.wantBooks[i],
							repo.books[i],
						)
					}

					if tt.nextID != repo.nextID {
						t.Fatalf(
							"expected nextID %d but got %d",
							tt.wantNextID,
							repo.nextID,
						)
					}

					for i := range tt.books {
						tt.books[i].ID++
					}

					for i := range tt.books {
						if tt.books[i] == repo.books[i] {
							t.Fatalf(
								"got unexpectedly changed book %v but expected but expected %v",
								repo.books[i],
								tt.wantBooks[i],
							)
						}
					}

				}

			}

		})
	}
}
