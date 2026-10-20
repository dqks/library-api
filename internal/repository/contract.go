package repository

type GetBooksPayload struct {
	Available *bool
}

type EditBookPayload struct {
	Title     *string
	Author    *string
	Year      *uint16
	Available *bool
}

type CreateBookPayload struct {
	Title     *string
	Author    *string
	Year      *uint16
	Available *bool
}
