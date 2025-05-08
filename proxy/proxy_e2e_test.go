package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

func startTestProxyServer() *httptest.Server {
	// Create a dummy target server to forward requests to
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Create the proxy server

	targetURL, err := url.Parse(target.URL)
	if err != nil {
		panic(err)
	}
	config := Config{
		Target:     "http://localhost:8080",
		ProxyPort:  8081,
		RateLimit:  1,
		BurstLimit: 1,
		CacheSize:  2,
	}
	rateLimiter := NewRateLimiter(config)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Wrap the proxy with middlewares
	handler := limitMiddleware(logMiddleware(proxy), rateLimiter)

	// Start the test server
	return httptest.NewServer(handler)
}

func TestRateLimiterLRUE2E(t *testing.T) {
	server := startTestProxyServer()
	defer server.Close()

	client := &http.Client{}

	// Simulate requests from different IPs
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	for _, ip := range ips {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set("X-Forwarded-For", ip)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed for IP %s: %v", ip, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK for IP %s, got %d", ip, resp.StatusCode)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	// Send another request from the first IP to ensure it's still allowed
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed for IP 192.168.1.1: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for IP 192.168.1.1, got %d", resp.StatusCode)
	}

	// Verify eviction: the second IP should now be evicted
	req, _ = http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.2")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Request failed for IP 192.168.1.2: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Error("Expected rate limiter to block IP 192.168.1.2 after eviction")
	}
}
