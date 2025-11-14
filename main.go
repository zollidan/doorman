package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zollidan/doorman/config"
	"github.com/zollidan/doorman/database"
	"github.com/zollidan/doorman/handlers"
)

func main() {

	cfg := config.New()

	db := database.Init(cfg)

	handlers := handlers.New(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		// r.Post("/login", handlers.LoginHandler)
		r.Post("/register",	handlers.RegisterHandler)
		// r.Get("/verify", handlers.VerifyHandler)
		
	})

	fmt.Printf("Server is running on %s", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}
