package main

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/handler"
	"library-api/internal/repository/postgres"

	"library-api/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)

	}
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	defer pool.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatal(err)
	}

	pgRepo := postgres.Create(pool)

	mux := http.NewServeMux()

	service := service.Create(pgRepo)

	if service == nil {
		fmt.Println("произошла ошибка при запуске сервера - не получилось создать сервис")
		return
	}

	mux.HandleFunc("POST /books", handler.CreateBook(service))
	mux.HandleFunc("GET /books", handler.GetBooks(service))
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID(service))
	mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID(service))
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID(service))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Println("Сервер запущен на порту 8080")
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println("Сервер закрылся штатно")
			} else {
				fmt.Println(err.Error())
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}
			return
		}
	}()

	<-shutDownSignal
	fmt.Println("Сигнал закрытия сервера получен, выполняется закрытие")
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutDownCtx); err != nil {
		fmt.Println("Произошла ошибка во время закрытия сервера")
	}
}
