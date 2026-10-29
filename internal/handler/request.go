package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"library-api/internal/apperrors"
	"net/http"
	"strconv"
)

func DecodeBody(requestBody io.ReadCloser, body interface{}, w http.ResponseWriter) error {
	decoder := json.NewDecoder(requestBody)
	decoder.DisallowUnknownFields()
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

func ValidatePathID(w http.ResponseWriter, pathID string) (int, error) {
	id, err := strconv.Atoi(pathID)
	if err != nil || id <= 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(apperrors.ErrIncorrectPathValueCode)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()}); err != nil {
			fmt.Println(err)
		}
		return -1, apperrors.ErrIncorrectPathValue
	}
	return id, nil
}
