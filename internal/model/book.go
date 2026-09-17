package model

import (
	"library-api/internal/apperrors"
	"strings"
	"time"
)

type Book struct {
	ID        int
	Title     string
	Author    string
	Year      uint16
	Available bool
}

type BookDTO struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      uint16 `json:"year"`
	Available bool   `json:"available"`
}

func (b *Book) ToDTO() BookDTO {
	return BookDTO{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		Year:      b.Year,
		Available: b.Available,
	}
}

func BookDomainListToDTO(books []Book) []*BookDTO {
	booksDTO := make([]*BookDTO, len(books))

	for i := range books {
		b := books[i].ToDTO()
		booksDTO[i] = &b
	}

	return booksDTO
}

func ValidateBookFields(title *string, author *string, year *uint16, available *bool) error {
	if author == nil && available == nil && title == nil && year == nil {
		return apperrors.ErrRequiredFields
	}

	if (author != nil && len(strings.TrimSpace(*author)) == 0) || (title != nil && len(strings.TrimSpace(*title)) == 0) {
		return apperrors.ErrInvalidValues
	}

	if year != nil && *year > uint16(time.Now().Year()) {
		return apperrors.ErrInvalidValues
	}
	return nil
}
