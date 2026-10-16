package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestLifecycleBook(t *testing.T) {
	repo := repository.Create(3, []model.Book{
		{
			ID:        1,
			Title:     "Peace and war",
			Author:    "Lev Tolstoy",
			Year:      2000,
			Available: false,
		},
		{
			ID:        2,
			Title:     "New title 2",
			Author:    "New Author 2",
			Year:      2000,
			Available: true,
		},
	})
	service := service.Create(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", handler.GetBooks(service))
	mux.HandleFunc("POST /books", handler.CreateBook(service))
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID(service))
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID(service))
	mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID(service))

	expectedHeaders := map[string][]string{
		"Content-Type": {"application/json"},
	}

	// Create book
	createBody := handler.CreateBookBody{
		Title:     new("title"),
		Author:    new("author"),
		Year:      new(uint16(2000)),
		Available: new(true),
	}
	ctx := context.WithoutCancel(context.Background())
	createRr := httptest.NewRecorder()
	var body bytes.Buffer
	encoder := json.NewEncoder(&body)
	err := encoder.Encode(createBody)
	if err != nil {
		t.Fatal(
			"failed to encode body",
		)
	}
	createReq := httptest.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/books",
		&body,
	)

	mux.ServeHTTP(createRr, createReq)

	if createRr.Code != 201 {
		t.Fatalf(
			"expected code to be 201 but got %d",
			createRr.Code,
		)
	}

	if !testutil.CompareHeaders(expectedHeaders, createRr.Header()) {
		t.Fatalf(
			"expected headers to be %v\nbut got %v",
			expectedHeaders,
			createRr.Header(),
		)
	}

	// Get book by ID
	var createdBook model.Book
	json.NewDecoder(createRr.Body).Decode(&createdBook)
	getBookRr := httptest.NewRecorder()
	getBookReq := httptest.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/books/"+strconv.Itoa(createdBook.ID),
		nil,
	)
	mux.ServeHTTP(getBookRr, getBookReq)

	var getBook model.Book
	decoder := json.NewDecoder(getBookRr.Body)
	err = decoder.Decode(&getBook)

	if err != nil {
		t.Fatal(
			"failed to decode body",
		)
	}

	if createdBook != getBook {
		t.Fatalf(
			"expected found by id %d book to be %v\nbut got %v",
			createdBook.ID,
			createdBook,
			getBook,
		)
	}

	if getBookRr.Code != 200 {
		t.Fatalf(
			"expected code to be 200 but got %d",
			getBookRr.Code,
		)
	}

	if !testutil.CompareHeaders(expectedHeaders, getBookRr.Header()) {
		t.Fatalf(
			"expected headers to be %v\nbut got %v",
			expectedHeaders,
			getBookRr.Header(),
		)
	}

	// Edit book
	editBody := handler.EditBookBody{
		Title: new("edited title"),
	}
	var editBodyJSON bytes.Buffer
	json.NewEncoder(&editBodyJSON).Encode(editBody)
	editRr := httptest.NewRecorder()
	editReq := httptest.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		"/books/"+strconv.Itoa(createdBook.ID),
		&editBodyJSON,
	)
	mux.ServeHTTP(editRr, editReq)
	var editedBook model.Book
	decoder = json.NewDecoder(editRr.Body)
	err = decoder.Decode(&editedBook)
	if err != nil {
		t.Fatal(
			"failed to decode body",
		)
	}

	createdBook.Title = "edited title"
	if editedBook != createdBook {
		t.Fatalf(
			"expected edited book to be %v\nbut got %v",
			createdBook, editedBook,
		)
	}

	if editRr.Code != 200 {
		t.Fatalf(
			"expected code to be 200 but got %d",
			editRr.Code,
		)
	}

	if !testutil.CompareHeaders(expectedHeaders, editRr.Header()) {
		t.Fatalf(
			"expected headers to be %v\nbut got %v",
			expectedHeaders,
			editRr.Header(),
		)
	}

	// Get books
	getBooksRr := httptest.NewRecorder()
	getBooksReq := httptest.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/books",
		nil,
	)
	mux.ServeHTTP(getBooksRr, getBooksReq)

	var getBooks []model.Book
	decoder = json.NewDecoder(getBooksRr.Body)
	err = decoder.Decode(&getBooks)
	if err != nil {
		t.Fatal(
			"failed to decode body",
		)
	}

	foundBook := false
	for i := range getBooks {
		if getBooks[i] == createdBook {
			foundBook = true
		}
	}

	if !foundBook {
		t.Fatalf(
			"expected to find book in list of all books %v\nbut didn't %v",
			createdBook,
			getBooks,
		)
	}

	if getBooksRr.Code != 200 {
		t.Fatalf(
			"expected code to be 200 but got %d",
			getBooksRr.Code,
		)
	}

	if !testutil.CompareHeaders(expectedHeaders, getBooksRr.Header()) {
		t.Fatalf(
			"expected headers to be %v\nbut got %v",
			expectedHeaders,
			getBooksRr.Header(),
		)
	}

	// Delete book
	delRr := httptest.NewRecorder()
	delReq := httptest.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		"/books/"+strconv.Itoa(createdBook.ID),
		nil,
	)
	mux.ServeHTTP(delRr, delReq)

	if delRr.Code != 200 {
		t.Fatalf(
			"expected code to be 200 but got %d",
			delRr.Code,
		)
	}

	if !testutil.CompareHeaders(expectedHeaders, delRr.Header()) {
		t.Fatalf(
			"expected headers to be %v\nbut got %v",
			expectedHeaders,
			delRr.Header(),
		)
	}

	// Get book by id
	getBookRr = httptest.NewRecorder()
	getBookReq = httptest.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/books/"+strconv.Itoa(createdBook.ID),
		nil,
	)
	mux.ServeHTTP(getBookRr, getBookReq)

	if getBookRr.Code != 404 {
		t.Fatalf(
			"expected code to be 404 but got %d",
			getBookRr.Code,
		)
	}

	if !testutil.CompareHeaders(expectedHeaders, getBookRr.Header()) {
		t.Fatalf(
			"expected headers to be %v\nbut got %v",
			expectedHeaders,
			getBookRr.Header(),
		)
	}
}
