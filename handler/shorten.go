package handler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type URLHandler struct {
	rdb *redis.Client
}

func NewURLHandler(rdb *redis.Client) *URLHandler {
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

	fmt.Fprintf(w, "http://localhost:8080/r/%s", shortURL)
}

func (h *URLHandler) RedirectShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	shortURL := chi.URLParam(r, "shortURL")
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
