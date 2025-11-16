package router

import (
	"github.com/Gylmynnn/golink/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/metrics", handlers.GetTopDomains)
	r.Post("/shorten", handlers.ShortenURL)
	r.Get("/{shortURL}", handlers.RedirectURL)
	return r
}
