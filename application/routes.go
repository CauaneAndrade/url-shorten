package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/CauaneAndrade/url-shorten/handler"
)

func loadRoutes(rdb *redis.Client) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Wrap loadShortenRoutes inside an anonymous function
	router.Route("/", func(r chi.Router) {
		loadShortenRoutes(r, rdb)
	})

	return router
}

// Adjust loadShortenRoutes to accept the Redis client
func loadShortenRoutes(router chi.Router, rdb *redis.Client) {
	urlHandler := handler.NewURLHandler(rdb) // Initialize URLHandler with Redis client

	router.Post("/shorten", urlHandler.GenerateShortURL)
	router.Get("/r/{shortURL}", urlHandler.RedirectShortURL)
}
