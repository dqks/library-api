package service_test

import (
	"library-api/internal/service"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		repo        service.BookRepository
		serviceNil  bool
		wantService service.BookService
	}{
		{
			name:        "success, service nil",
			repo:        nil,
			serviceNil:  true,
			wantService: service.BookService{},
		},
		{
			name: "success",
			repo: &MockRepo{},
			wantService: service.BookService{
				Repo: &MockRepo{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := service.CreateBookService(tt.repo)

			if tt.serviceNil {
				if service != nil {
					t.Fatalf("expected service to be %v but got %v", tt.wantService, service)
				}
			} else {
				if !reflect.DeepEqual(service, &tt.wantService) {
					t.Fatalf("expected service to be %#v but got %#v", &tt.wantService, service)
				}
			}
		})
	}
}
