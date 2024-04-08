package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

// type URLHandler struct {
// 	rdb *redis.Client
// }

// func NewURLHandler(rdb *redis.Client) *URLHandler {
// 	return &URLHandler{rdb: rdb}
// }

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
}

type URLHandler struct {
	rdb RedisClient
}

func NewURLHandler(rdb RedisClient) *URLHandler {
	return &URLHandler{rdb: rdb}
}

// GenerateShortURL creates a short URL and stores it in the handler's map
func (h *URLHandler) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Validate URL
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURL := generateRandomString(5)
	err = h.rdb.Set(ctx, shortURL, originalURL, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save short URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"short_url": "http://localhost:8080/r/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *URLHandler) IncrementAccessCount(ctx context.Context, shortURL string) error {
	countKey := "count:" + shortURL
	if _, err := s.rdb.Incr(ctx, countKey).Result(); err != nil {
		return fmt.Errorf("error incrementing access count: %w", err)
	}
	return nil
}

func (h *URLHandler) RedirectShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	shortURL := chi.URLParam(r, "shortURL")

	// Increment the access count for this short URL
	if err := h.IncrementAccessCount(ctx, shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	originalURL, err := h.rdb.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving short URL", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
