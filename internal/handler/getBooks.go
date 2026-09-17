package handler

import (
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	w.Header().Set("Content-Type", "application/json")
	availParam := false
	encoder := json.NewEncoder(w)
	var books []*model.Book

	if availQuery := r.URL.Query().Get("available"); availQuery != "" {
		var err error
		availParam, err = strconv.ParseBool(availQuery)
		if err != nil {
			w.WriteHeader(400)
			err := encoder.Encode(apperrors.BaseError{Error: err.Error()})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			return
		}
		books = service.GetBooks(context, service.GetBooksQueryParams{Available: &availParam})
	} else {
		books = service.GetBooks(context, service.GetBooksQueryParams{Available: nil})
	}

	if books == nil {
		w.WriteHeader(500)
		encoder.Encode(apperrors.BaseError{Error: apperrors.ErrInternal.Error()})
		return
	}

	err := encoder.Encode(model.BookDomainListToDTO(books))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
