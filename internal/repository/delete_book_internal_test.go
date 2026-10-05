package repository

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestDeleteBookByIDConcurrentCancel(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "success",
			id:      1,
			wantErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}
			var err error
			testChan := make(chan bool)
			waitChan := make(chan bool)

			repo.beforeLock = func() {
				testChan <- true
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				<-waitChan

				go func() {
					defer wg.Done()
					err = repo.DeleteBookByID(ctx, tt.id)
				}()

				<-testChan
				waitChan <- true
			}()

			go func() {
				defer wg.Done()
				repo.mutex.Lock()
				waitChan <- true
				<-waitChan
				cancel()
				repo.mutex.Unlock()
			}()

			wg.Wait()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v but expected %v", err, tt.wantErr)
			}

			ctxGet := context.WithoutCancel(context.Background())
			books, err := repo.GetBooks(ctxGet, GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			foundBook := false
			for _, b := range books {
				if b.ID == tt.id {
					foundBook = true
				}
			}

			if !foundBook {
				t.Fatalf("book by id %d was deleted", tt.id)
			}
		})
	}
}
