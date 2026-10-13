package tests

import (
	"context"
	"library-api/internal/handler"
	"library-api/internal/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		cancel          bool
		deadlineContext bool
		serviceErr      error
		serviceBook     model.Book
		wantId          string
		wantBody        string
		wantHeadrs      http.Header
		wantCode        int
	}{
		{
			name: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
		var ctx context.Context
		var cancel context.CancelFunc

		if tt.deadlineContext {
			ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
			defer cancel()
		} else {
			ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
			defer cancel()
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://localhost:8080/books/"+strconv.Itoa(tt.id),
			nil,
		)
		req.SetPathValue("id", strconv.Itoa(tt.id))

		s := ServiceMock{
			err:  tt.serviceErr,
			book: tt.serviceBook,
		}

		handler.GetBookByID(&s)(rr, req)

	}
}
