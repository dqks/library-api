package handler

import (
	"encoding/json"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	books := service.GetBooks(context)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(model.BookDomainListToDTO(books))
}
