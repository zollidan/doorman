package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zollidan/doorman/config"
	"github.com/zollidan/doorman/database"
)

func main() {

	cfg := config.New()
	_ = database.Init(cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})
	})

	fmt.Printf("Server is running on %s", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}
