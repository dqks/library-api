package service

type bookService struct {
	repo BookRepository
}

func CreateBookService(repo BookRepository) *bookService {
	if repo == nil {
		return nil
	}

	return &bookService{repo: repo}
}
