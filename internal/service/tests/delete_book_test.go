package tests

import (
	"context"
	"library-api/internal/service"
	"testing"
)

func TestDeleteBookByID(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		cancel  bool
		wantErr error
	}{
		{
			name:    "success",
			id:      1,
			cancel:  false,
			wantErr: nil,
		},
		{
			name:    "canceled",
			id:      1,
			cancel:  true,
			wantErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
				cancel()
			}

			s := service.Create(&MockRepo{})

			err := s.DeleteBookByID(ctx, tt.id)

			if err != tt.wantErr {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}
		})
	}
}
