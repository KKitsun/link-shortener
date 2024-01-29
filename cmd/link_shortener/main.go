package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/KKitsun/link-shortener/internal/config"
	"github.com/KKitsun/link-shortener/internal/handlers/getFull"
	"github.com/KKitsun/link-shortener/internal/handlers/shorten"
	"github.com/KKitsun/link-shortener/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	storage, err := postgres.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to init storage: %s", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/url", func(r chi.Router) {
		r.Post("/", shorten.Shorten(storage))
		r.Get("/{alias}", getFull.GetFull(storage))
	})

	if err := http.ListenAndServe(":"+cfg.HTTPServer.Port, r); err != nil {
		log.Fatalf("failed to start server")
	}
}
