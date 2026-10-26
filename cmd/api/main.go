package main

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/database"
	"library-api/internal/handler"
	"library-api/internal/repository/postgres"
	"log"

	"library-api/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	dsn, err := database.GetDsn("../../.env", "DATABASE_URL")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := database.CreatePool(ctx, dsn)

	if err != nil {
		log.Fatal(err)
	}

	pgRepo := postgres.Create(pool)

	mux := http.NewServeMux()

	bookService := service.CreateBookService(pgRepo)

	if bookService == nil {
		fmt.Println("произошла ошибка при запуске сервера - не получилось создать сервис")
		return
	}

	userService := service.CreateUserService(pgRepo)

	if userService == nil {
		fmt.Println("произошла ошибка при запуске сервера - не получилось создать сервис")
		return
	}

	mux.HandleFunc("POST /books", handler.CreateBook(bookService))
	mux.HandleFunc("GET /books", handler.GetBooks(bookService))
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID(bookService))
	mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID(bookService))
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID(bookService))
	mux.HandleFunc("POST /users", handler.CreateUser(userService))

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
