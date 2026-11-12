package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type readyResponse struct {
	Error error `json:"error"`
}

func Ready(pool *pgxpool.Pool) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		w.Header().Add("Content-Type", "application/json")
		err := pool.Ping(ctx)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(readyResponse{Error: err}); err != nil {
			slog.Error(
				"ready request error",
				"err", err,
			)
		}
	}
}
