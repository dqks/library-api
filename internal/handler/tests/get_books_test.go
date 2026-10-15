package tests

import (
	"context"
	"errors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetBooks(t *testing.T) {
	tests := []struct {
		name                   string
		availableParam         string
		cancel                 bool
		deadlineContext        bool
		serviceErr             error
		serviceBooks           []model.Book
		wantCode               int
		wantServiceCallCount   int
		wantHeaders            http.Header
		wantBody               string
		wantParamsAvailable    bool
		wantParamsAvailableNil bool
		wantBooks              []model.Book
	}{
		{
			name:           "success, available true",
			availableParam: "?available=true",
			serviceErr:     nil,
			serviceBooks: []model.Book{
				{
					ID:        1,
					Title:     "title",
					Author:    "title",
					Year:      2000,
					Available: true,
				},
			},
			wantCode:             200,
			wantServiceCallCount: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody:               `[{"id":1,"title":"title","author":"title","year":2000,"available":true}]`,
			wantParamsAvailable:    true,
			wantParamsAvailableNil: false,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "title",
					Author:    "title",
					Year:      2000,
					Available: true,
				},
			},
		},
		{
			name:           "success, available false",
			availableParam: "?available=false",
			serviceErr:     nil,
			serviceBooks: []model.Book{
				{
					ID:        1,
					Title:     "title",
					Author:    "title",
					Year:      2000,
					Available: false,
				},
			},
			wantCode:             200,
			wantServiceCallCount: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody:               `[{"id":1,"title":"title","author":"title","year":2000,"available":false}]`,
			wantParamsAvailable:    true,
			wantParamsAvailableNil: false,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "title",
					Author:    "title",
					Year:      2000,
					Available: false,
				},
			},
		},
		{
			name:           "success, no params",
			availableParam: "",
			serviceErr:     nil,
			serviceBooks: []model.Book{
				{
					ID:        1,
					Title:     "title",
					Author:    "title",
					Year:      2000,
					Available: true,
				},
			},
			wantCode:             200,
			wantServiceCallCount: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody:               `[{"id":1,"title":"title","author":"title","year":2000,"available":true}]`,
			wantParamsAvailable:    true,
			wantParamsAvailableNil: false,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "title",
					Author:    "title",
					Year:      2000,
					Available: true,
				},
			},
		},
		{
			name:                 "cancelled in handler",
			availableParam:       "",
			serviceErr:           nil,
			cancel:               true,
			wantCode:             200,
			wantServiceCallCount: 0,
			wantHeaders:          map[string][]string{},
		},
		{
			name:                 "deadline exceeded in handler",
			availableParam:       "",
			serviceErr:           nil,
			deadlineContext:      true,
			wantCode:             504,
			wantServiceCallCount: 0,
			wantHeaders:          map[string][]string{},
			wantBody:             ``,
		},
		{
			name:                 "cancelled in service",
			availableParam:       "",
			serviceErr:           context.Canceled,
			wantCode:             200,
			wantServiceCallCount: 1,
			wantHeaders:          map[string][]string{},
		},
		{
			name:                 "deadline exceeded in service",
			availableParam:       "",
			serviceErr:           context.DeadlineExceeded,
			wantCode:             504,
			wantServiceCallCount: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"context deadline exceeded"}`,
		},
		{
			name:                 "bad params",
			availableParam:       "?available=truedfkgdyr",
			serviceErr:           context.DeadlineExceeded,
			wantCode:             400,
			wantServiceCallCount: 0,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"strconv.ParseBool: parsing \"truedfkgdyr\": invalid syntax"}`,
		},
		{
			name:       "unknown service err",
			serviceErr: errors.New("unknown"),
			wantBody:   `{"error":"unknown"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             500,
			wantServiceCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			var cancel context.CancelFunc

			if tt.deadlineContext {
				ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()
			} else {
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
			}

			if tt.cancel {
				cancel()
			}

			rr := httptest.NewRecorder()

			var url = "http://localhost:8080/books" + tt.availableParam

			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodGet,
				url,
				nil,
			)

			s := ServiceMock{
				books: tt.serviceBooks,
				err:   tt.serviceErr,
			}

			handler.GetBooks(&s)(rr, req)

			if ctx != s.ctx && tt.wantServiceCallCount >= 1 {
				t.Error(
					"expected service's context to be the same as handler's one",
				)
			}

			if tt.wantCode != rr.Code {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if !testutil.CompareHeaders(tt.wantHeaders, rr.Header()) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

			if tt.wantServiceCallCount != s.callCount {
				t.Errorf(
					"expected service to be called %d times but got %d",
					tt.wantServiceCallCount,
					s.callCount,
				)
			}

			if strings.TrimSpace(tt.wantBody) != strings.TrimSpace(rr.Body.String()) {
				t.Errorf(
					"expected body to be %v\nbut got %v",
					tt.wantBody,
					rr.Body.String(),
				)
			}

			if tt.wantServiceCallCount >= 1 {
				if len(tt.wantBooks) != len(s.books) {
					t.Errorf(
						"expected to use books from repository\nexpected books %v\n actual books %v",
						tt.wantBooks,
						s.books,
					)
				}
				for i := range tt.wantBooks {
					if tt.wantBooks[i] != s.books[i] {
						t.Errorf(
							"expected to use books from repository\nexpected books %v\n actual books %v",
							tt.wantBooks,
							s.books,
						)
						break
					}
				}
			}

			if (tt.availableParam == "?available=true" || tt.availableParam == "?available=false" || tt.availableParam == "") &&
				tt.wantServiceCallCount >= 1 {
				var expectedParams service.GetBooksQueryParams

				switch tt.availableParam {
				case "?available=true":
					trueParam := true
					expectedParams.Available = &trueParam
				case "?available=false":
					trueParam := false
					expectedParams.Available = &trueParam
				}

				if expectedParams.Available != nil && s.getParams.Available != nil {
					if *expectedParams.Available != *s.getParams.Available {
						t.Errorf(
							"expected request params to be %v\nbut got %v",
							expectedParams,
							s.getParams,
						)
					}
				} else if expectedParams.Available == nil && s.getParams.Available != nil {
					t.Errorf(
						"expected request params to be %v\nbut got %v",
						expectedParams,
						s.getParams,
					)
				} else if expectedParams.Available != nil && s.getParams.Available == nil {
					t.Errorf(
						"expected request params to be %v\nbut got %v",
						expectedParams,
						s.getParams,
					)
				}
			}
		})
	}
}
