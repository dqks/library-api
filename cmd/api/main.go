package main

import (
	"library-api/internal/handler"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /books", handler.CreateBook)
	mux.HandleFunc("GET /books", handler.GetBooks)
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID)
	mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID)
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID)

	http.ListenAndServe(":8080", mux)
}
