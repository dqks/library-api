package main

import (
	"context"
	"errors"
	"library-api/config"
	"library-api/internal/database"
	"library-api/internal/handler"
	"library-api/internal/repository/postgres"
	"log/slog"
	"strconv"

	"library-api/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slogHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	config, err := config.Get(true)

	if err != nil {
		slog.Error(
			"failed to get config",
			"err", err,
		)
		return
	}

	ctx := context.Background()

	pool, err := database.CreatePool(ctx, config.DSN)

	if err != nil {
		slog.Error(
			"failed to create pool",
			"err", err,
		)
		return
	}

	pgRepo := postgres.Create(pool)

	mux := http.NewServeMux()

	bookService := service.CreateBookService(pgRepo)

	if bookService == nil {
		slog.Error(
			"failed to create book service",
		)
		return
	}

	userService := service.CreateUserService(pgRepo)

	if userService == nil {
		slog.Error(
			"failed to create user service",
		)
		return
	}

	loanService := service.CreateLoanService(pgRepo)

	if loanService == nil {
		slog.Error(
			"failed to create loan service",
		)
		return
	}

	mux.HandleFunc("POST /books", handler.CreateBook(bookService))
	mux.HandleFunc("GET /books", handler.GetBooks(bookService))
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID(bookService))
	mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID(bookService))
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID(bookService))
	mux.HandleFunc("POST /users", handler.CreateUser(userService))
	mux.HandleFunc("GET /users/{id}", handler.GetUser(userService))
	mux.HandleFunc("POST /loans/{bookId}/{userId}", handler.LoanBook(loanService))
	mux.HandleFunc("PATCH /loans/{loanId}", handler.ReturnBook(loanService))
	mux.HandleFunc("GET /loans/{id}", handler.GetLoan(loanService))
	mux.HandleFunc("GET /users/{userId}/loans", handler.GetUserLoans(loanService))

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: mux,
	}

	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		slog.Info(
			"server stared",
			"port", strconv.Itoa(config.Port),
		)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Info(
				"server failed to start",
				"err", err,
			)
			return
		}
	}()
	<-shutDownSignal
	slog.Info(
		"server closing",
	)
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutDownCtx); err != nil {
		slog.Info(
			"server failed to close",
			"err", err,
		)
	}
}
