package main

import (
	"log"
	"net/http"

	"github.com/Gylmynnn/golink/router"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", router.SetupRouter()) 

	log.Println("Starting server on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
