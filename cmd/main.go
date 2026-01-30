package main

import (
	"context"
	"finance-app/db"
	"finance-app/internal/auth"
	"finance-app/internal/repository"
	middlewares "finance-app/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// LOAD ENV
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		_ = godotenv.Load()
	}

	// DB CONNECTION
	makeConnection := &db.Connection{}
	pgxPool, err := makeConnection.Connect(context.Background())
	if err != nil {
		log.Println("failed to connect database:", err)
	}
	defer pgxPool.Close()

	// MODULES
	repository := repository.New(pgxPool)
	authService := auth.NewService(repository)

	// ROUTES
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("health is ok")
		w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middlewares.CorsMiddleware)
		auth.NewHandler(r, authService)
	})

	// PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("Server running production, on port", port)
	http.ListenAndServe(":"+port, r)
}
