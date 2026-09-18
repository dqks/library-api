package handler

import "library-api/internal/model"

type BookDTO struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      uint16 `json:"year"`
	Available bool   `json:"available"`
}

func BookModelToDTO(b *model.Book) BookDTO {
	return BookDTO{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		Year:      b.Year,
		Available: b.Available,
	}
}

func BookDomainListToDTO(books []model.Book) []*BookDTO {
	booksDTO := make([]*BookDTO, len(books))

	for i := range books {
		b := BookModelToDTO(&books[i])
		booksDTO[i] = &b
	}

	return booksDTO
}
