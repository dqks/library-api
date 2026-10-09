package tests

import (
	"encoding/json"
	"fmt"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		getInService   bool
		serviceErr     error
		serviceBook    model.Book
		wantCode       int
		wantHeader     map[string][]string
		wantResultBody string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(
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

			// fmt.Println(rr)

			if rr.Code != tt.wantCode {
				t.Fatalf("expected code %d but got %d", tt.wantCode, rr.Code)
			}

			if reflect.DeepEqual(tt.wantHeader, rr.Header()) {
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

			fmt.Println(tt.wantResultBody)
			fmt.Println(rr.Body.String())

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
