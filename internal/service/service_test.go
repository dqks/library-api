package service

import (
	"library-api/internal/service/tests"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		repo        Repository
		serviceNil  bool
		wantService BookService
	}{
		{
			name:        "success, service nil",
			repo:        nil,
			serviceNil:  true,
			wantService: BookService{},
		},
		{
			name: "success",
			repo: &tests.MockRepo{},
			wantService: BookService{
				repo: &tests.MockRepo{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := Create(tt.repo)

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
