package model

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
