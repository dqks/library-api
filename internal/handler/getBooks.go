package handler

import (
	"library-api/internal/service"
	"net/http"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	service.GetBooks(context)
}
