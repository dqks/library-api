package repository

import (
	"context"
	"fmt"
	"library-api/internal/model"
	"sync"
	"testing"
)

func TestGetBooksEdited(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		author    string
		year      uint16
		available bool
		wantBooks []model.Book
	}{
		{
			name:      "success",
			title:     "New title",
			author:    "New Author",
			year:      2000,
			available: true,
			wantBooks: []model.Book{
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
				{
					ID:        3,
					Title:     "New title",
					Author:    "New Author",
					Year:      uint16(2000),
					Available: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			waitChan := make(chan bool)
			repo.afterLock = func() {
				waitChan <- true
			}
			var books []model.Book
			var wg sync.WaitGroup
			wg.Add(2)
			ctx := context.WithoutCancel(context.Background())

			go func() {
				defer wg.Done()
				<-waitChan
				books, _ = repo.GetBooks(ctx, GetBooksPayload{})
			}()

			go func() {
				defer wg.Done()
				repo.CreateBook(ctx, CreateBookPayload{
					Title:     &tt.title,
					Author:    &tt.author,
					Year:      &tt.year,
					Available: &tt.available,
				})
			}()

			wg.Wait()

			if len(books) != len(tt.wantBooks) {
				t.Fatalf(
					"expected len %d but got %d",
					len(tt.wantBooks),
					len(books),
				)
			}

			for i := range books {
				if books[i] != tt.wantBooks[i] {
					t.Fatalf("got books %v \n but expected %v", books, tt.wantBooks)
				}
			}

		})
	}
}

func TestGetBooksConcurrentCancel(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			var err error
			var books []model.Book
			waitChan := make(chan bool, 1)
			testChan := make(chan bool, 1)
			repo.beforeLock = func() {
				testChan <- true
			}
			ctx, cancel := context.WithCancel(context.Background())
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				<-waitChan
				go func() {
					defer wg.Done()
					books, err = repo.GetBooks(ctx, GetBooksPayload{})
				}()

				<-testChan
				waitChan <- true
			}()

			go func() {
				defer wg.Done()
				repo.mutex.RLock()
				waitChan <- true
				<-waitChan
				cancel()
				repo.mutex.RUnlock()
			}()

			wg.Wait()

			fmt.Println(books)
			fmt.Println(err)
		})
	}
}
