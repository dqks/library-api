package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"testing"
	"time"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name             string
		title            string
		wantTitleNil     bool
		author           string
		wantAuthorNil    bool
		year             uint16
		wantYearNil      bool
		available        bool
		wantAvailableNil bool
		wantCancel       bool
		wantBook         model.Book
		wantErr          error
	}{
		{
			name:       "success",
			title:      "New Title",
			author:     "New Author",
			year:       uint16(2000),
			available:  true,
			wantCancel: false,
			wantBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(2000),
				Available: true,
			},
			wantErr: nil,
		},
		{
			name:         "err, nil title",
			title:        "",
			author:       "New Author",
			year:         uint16(2000),
			available:    true,
			wantCancel:   false,
			wantBook:     model.Book{},
			wantErr:      apperrors.ErrRequiredFields,
			wantTitleNil: true,
		},
		{
			name:          "err, nil author",
			title:         "New Title",
			author:        "",
			year:          uint16(2000),
			available:     true,
			wantCancel:    false,
			wantBook:      model.Book{},
			wantErr:       apperrors.ErrRequiredFields,
			wantAuthorNil: true,
		},
		{
			name:        "err, nil year",
			title:       "New Title",
			author:      "",
			year:        uint16(0),
			available:   true,
			wantCancel:  false,
			wantBook:    model.Book{},
			wantErr:     apperrors.ErrRequiredFields,
			wantYearNil: true,
		},
		{
			name:             "err, nil available",
			title:            "New Title",
			author:           "",
			year:             uint16(0),
			available:        true,
			wantCancel:       false,
			wantBook:         model.Book{},
			wantErr:          apperrors.ErrRequiredFields,
			wantAvailableNil: true,
		},
		{
			name:       "err, empty string title",
			title:      "",
			author:     "New Author",
			year:       uint16(0),
			available:  true,
			wantCancel: false,
			wantBook:   model.Book{},
			wantErr:    apperrors.ErrRequiredFields,
		},
		{
			name:       "err, empty string author",
			title:      "",
			author:     "New Author",
			year:       uint16(0),
			available:  true,
			wantCancel: false,
			wantBook:   model.Book{},
			wantErr:    apperrors.ErrRequiredFields,
		},
		{
			name:       "err, year more than now",
			title:      "New Title",
			author:     "New Author",
			year:       uint16(time.Now().Year() + 1),
			available:  true,
			wantCancel: false,
			wantBook:   model.Book{},
			wantErr:    apperrors.ErrRequiredFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.wantCancel {
				cancel()
			}

			// book, err :=

		})
	}
}
