package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CauaneAndrade/url-shorten/handler"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd) // Ensure this matches your Redis client's Set method return type
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd) // Match Get method return type
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd) // Match Incr method return type
}

func TestGenerateShortURL(t *testing.T) {
	// It verifies the response status code and the body to ensure
	// they match expected values when a valid URL is provided.
	mockRedis := new(MockRedisClient)
	urlHandler := handler.NewURLHandler(mockRedis)

	// Set up expectations for Redis client
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusResult("", nil))

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", "/shorten?url=http://test.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(urlHandler.GenerateShortURL)

	// Serve the HTTP request to our handler
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := `{"short_url":"http://localhost:8080/r/`
	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want prefix %v", rr.Body.String(), expected)
	}

	mockRedis.AssertExpectations(t)
}
