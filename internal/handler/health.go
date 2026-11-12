package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(healthResponse{Status: "ok"}); err != nil {
		slog.Error(
			"health request error",
			"err", err,
		)
	}
}
