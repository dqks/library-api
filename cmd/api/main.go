package main

import (
	"context"
	"fmt"
	"library-api/internal/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /books", handler.CreateBook)
	mux.HandleFunc("GET /books", handler.GetBooks)
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID)
	mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID)
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Println("Сервер запущен на порту 8080")
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
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
