package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"net/http"
)

func ProcessServiceErr(w http.ResponseWriter, err error) {
	code := apperrors.CheckErrors(
		[]error{
			context.Canceled,
			context.DeadlineExceeded,
			apperrors.ErrInvalidValues,
		},
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
