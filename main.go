package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/piotsik/moviesiec-go/pkg/handler"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	handler, err := handler.Init()
	if err != nil {
		panic(err)
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/", handler.Hello)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", handler.UserGetAll)
			r.Get("/{id}", handler.UserGetByID)
			r.Post("/", handler.UserAdd)
			r.Put("/{id}", handler.UserUpdate)
			r.Delete("/{id}", handler.UserDelete)
			r.Get("/{id}/ratings", handler.UserGetByIDRatingsAll)

			r.Route("/login", func(r chi.Router) {
				r.Post("/", handler.UserAuthenticate)
			})
		})

		r.Route("/movies", func(r chi.Router) {
			r.Get("/", handler.MovieGetAll)
			r.Get("/{id}", handler.MovieGetByID)
			r.Post("/", handler.MovieAdd)
			r.Delete("/{id}", handler.MovieDelete)

			r.Route("/rate", func(r chi.Router) {
				r.Post("/", handler.MovieRate)
			})
		})
	})

	http.ListenAndServe(":8080", r)
}
