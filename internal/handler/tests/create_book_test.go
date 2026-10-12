package tests

import (
	"context"
	"encoding/json"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name            string
		body            string
		getInService    bool
		serviceErr      error
		serviceBook     model.Book
		cancel          bool
		contextDeadline bool
		wantCode        int
		wantHeader      http.Header
		wantResultBody  string
	}{
		{
			name:         "success",
			body:         `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			getInService: true,
			wantCode:     201,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr: nil,
			serviceBook: model.Book{
				ID:        1,
				Title:     "New title",
				Author:    "New Author",
				Year:      2003,
				Available: false,
			},
			wantResultBody: `{"id":1,"title":"New title","author":"New Author","year":2003,"available":false}`,
		},
		{
			name:           "canceled",
			body:           `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			getInService:   false,
			wantCode:       200,
			wantHeader:     map[string][]string{},
			cancel:         true,
			serviceErr:     nil,
			serviceBook:    model.Book{},
			wantResultBody: ``,
		},
		{
			name:         "err, body invalid values",
			body:         `{"title":"New title","author":"New Author","year":2077,"available":false}`,
			getInService: true,
			wantCode:     400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:     apperrors.ErrInvalidValues,
			serviceBook:    model.Book{},
			wantResultBody: `{"error":"` + apperrors.ErrInvalidValues.Error() + `"}`,
		},
		{
			name:         "err, body absence of required field",
			body:         `{"title":"New title","author":"New Author","available":false}`,
			getInService: true,
			wantCode:     400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:     apperrors.ErrRequiredFields,
			serviceBook:    model.Book{},
			wantResultBody: `{"error":"` + apperrors.ErrRequiredFields.Error() + `"}`,
		},
		{
			name:         "err, incorrect json",
			body:         `{"title":"New title","author":"New Author","available":`,
			getInService: false,
			wantCode:     400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:     nil,
			serviceBook:    model.Book{},
			wantResultBody: `{"error":"unexpected EOF"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"http://localhost:8080/books",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			s := ServiceMock{
				err:  tt.serviceErr,
				book: tt.serviceBook,
			}

			if tt.cancel {
				cancel()
			}

			handler.CreateBook(&s)(rr, req)

			if rr.Code != tt.wantCode {
				t.Fatalf("expected code %d but got %d", tt.wantCode, rr.Code)
			}

			if !compareHeaders(tt.wantHeader, rr.Header()) {
				t.Fatalf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeader,
					rr.Header(),
				)
			}

			if tt.getInService && s.gotIntoService {
				var expectedInput service.CreateBookInput
				decoder := json.NewDecoder(strings.NewReader(tt.body))
				decoder.Decode(&expectedInput)
				if !reflect.DeepEqual(expectedInput, s.createInput) {
					t.Fatalf(
						"expected to call service with input %v\nbut got %v",
						expectedInput,
						s.createInput,
					)
				}
			} else if !tt.getInService && s.gotIntoService {
				t.Fatal("expected not to get into service but got")
			} else if tt.getInService && !s.gotIntoService {
				t.Fatal("expected to get into service but didn't")
			}

			if strings.TrimSpace(tt.wantResultBody) != strings.TrimSpace(rr.Body.String()) {
				t.Fatalf(
					"expected result body %s\nbut got %s",
					tt.wantResultBody,
					rr.Body.String(),
				)
			}
		})
	}
}

func TestCreateBookDeadlineExceeded(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		getInService   bool
		serviceErr     error
		serviceBook    model.Book
		wantCode       int
		wantHeader     http.Header
		wantResultBody string
	}{
		{
			name:           "err, deadline exceeded",
			body:           `{"title":"New title","author":"New Author","available":false}`,
			getInService:   false,
			wantCode:       200,
			wantHeader:     map[string][]string{},
			serviceErr:     nil,
			serviceBook:    model.Book{},
			wantResultBody: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Time{})
			defer cancel()

			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"http://localhost:8080/books",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			s := ServiceMock{
				err:  tt.serviceErr,
				book: tt.serviceBook,
			}

			handler.CreateBook(&s)(rr, req)

			if rr.Code != tt.wantCode {
				t.Fatalf("expected code %d but got %d", tt.wantCode, rr.Code)
			}

			if !compareHeaders(tt.wantHeader, rr.Header()) {
				t.Fatalf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeader,
					rr.Header(),
				)
			}

			if tt.getInService && s.gotIntoService {
				var expectedInput service.CreateBookInput
				decoder := json.NewDecoder(strings.NewReader(tt.body))
				decoder.Decode(&expectedInput)
				if !reflect.DeepEqual(expectedInput, s.createInput) {
					t.Fatalf(
						"expected to call service with input %v\nbut got %v",
						expectedInput,
						s.createInput,
					)
				}
			} else if !tt.getInService && s.gotIntoService {
				t.Fatal("expected not to get into service but got")
			} else if tt.getInService && !s.gotIntoService {
				t.Fatal("expected to get into service but didn't")
			}

			if strings.TrimSpace(tt.wantResultBody) != strings.TrimSpace(rr.Body.String()) {
				t.Fatalf(
					"expected result body %s\nbut got %s",
					tt.wantResultBody,
					rr.Body.String(),
				)
			}
		})
	}
}
