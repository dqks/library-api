package handler

import (
	"encoding/json"
	"io"
	"library-api/internal/apperrors"
	"net/http"
)

func DecodeBody(requestBody io.ReadCloser, body interface{}, w http.ResponseWriter) error {
	decoder := json.NewDecoder(requestBody)
	if err := decoder.Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apperrors.ErrInvalidValuesCode)
		encoder := json.NewEncoder(w)
		if baseErr := encoder.Encode(apperrors.BaseError{Error: err.Error()}); baseErr != nil {
			return baseErr
		}
		return err
	}
	return nil
}
