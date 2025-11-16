package router

import (
	"github.com/Gylmynnn/golink/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.IndexPage)
	r.Get("/metrics", handlers.GetTopDomains)
	r.Post("/shorten", handlers.ShortenURL)
	r.Get("/{shortURL}", handlers.RedirectURL)

	r.Post("/shorten-html", handlers.ShortenURLHTML)
	r.Get("/metrics-html", handlers.GetTopDomainsHTML)

	return r
}
