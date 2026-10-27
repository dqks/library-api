package handler

import (
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"net/http"
)

func ProcessServiceErr(w http.ResponseWriter, err error, errors []error) {
	code := apperrors.CheckErrors(
		errors,
		err,
	)

	if code == -1 {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
		fmt.Println(apperrors.ErrInternal.Error())
	}
}
