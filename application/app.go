package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	app := &App{
		router: loadRoutes(rdb), // Pass the Redis client to loadRoutes
		rdb:    rdb,
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	fmt.Println("Starting server")

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
