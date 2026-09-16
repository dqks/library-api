package handler

import (
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	w.Header().Set("Content-Type", "application/json")
	availParam := false
	encoder := json.NewEncoder(w)

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
	}

	books := service.GetBooks(context, repository.GetBooksQueryParams{Available: availParam})

	if books == nil {
		w.WriteHeader(500)
		encoder.Encode(apperrors.BaseError{Error: apperrors.ErrInternal.Error()})
	}

	err := encoder.Encode(model.BookDomainListToDTO(books))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
