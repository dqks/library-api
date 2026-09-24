package main

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	repo := repository.CreateRepo(
		3,
		[]model.Book{
			{
				ID:        1,
				Title:     "1984",
				Author:    "George Orwell",
				Year:      uint16(1956),
				Available: true,
			},
			{
				ID:        2,
				Title:     "Kallocain",
				Author:    "Karin Boye",
				Year:      uint16(1937),
				Available: false,
			},
		},
	)

	if repo == nil {
		fmt.Println("произошла ошибка при запуске сервера - не получилсь создать репозиторий")
		return
	}

	service := service.CreateBookService(repo)

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
