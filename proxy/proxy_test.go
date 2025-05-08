package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimitMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	config := setupTestConfig()
	rateLimiter := NewRateLimiter(config)
	handler := limitMiddleware(testHandler, rateLimiter)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}
	for i := range 5 {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "127.0.0.1:12345" // Simulate client IP

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if i == 0 {
			// First request should be allowed
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d\n", resp.StatusCode)
			}
		} else {
			// Subsequent requests should be rate-limited
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Expected status 429, got %d\n", resp.StatusCode)
			}
		}

		// Wait a short time
		time.Sleep(100 * time.Millisecond)
	}
}
