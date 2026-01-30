package main

import (
	"context"
	"finance-app/db"
	"finance-app/internal/auth"
	"finance-app/internal/repository"
	middlewares "finance-app/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	makeConnection := &db.Connection{}
	pgxPool, err := makeConnection.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	defer pgxPool.Close()

	repository := repository.New(pgxPool)
	authService := auth.NewService(repository)

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(middlewares.CorsMiddleware)
		auth.NewHandler(r, authService)
	})

	http.ListenAndServe(":8000", r)
}
