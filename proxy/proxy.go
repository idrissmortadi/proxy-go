package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/time/rate"
)

type Config struct {
	Target     string `yaml:"target"`
	ProxyPort  int    `yaml:"proxy_port"`
	RateLimit  int    `yaml:"rate_limit"`
	BurstLimit int    `yaml:"burst_limit"`
	CacheSize  int    `yaml:"cache_size"`
}

type RateLimiter struct {
	limitersCache *lru.Cache
	mu            sync.Mutex
	config        Config
}

func NewRateLimiter(config Config) *RateLimiter {
	limitersLRU, err := lru.New(config.CacheSize)
	if err != nil {
		panic(err)
	}
	return &RateLimiter{
		limitersCache: limitersLRU,
		config:        config,
	}
}

func (rl *RateLimiter) GetLimiter(clientIP string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if a limiter exists for the client IP
	if rl.limitersCache.Contains(clientIP) {
		limiter, _ := rl.limitersCache.Get(clientIP)
		return limiter.(*rate.Limiter)
	}

	// Create a new limiter if none exists
	limiter := rate.NewLimiter(rate.Limit(rl.config.RateLimit), rl.config.BurstLimit)
	rl.limitersCache.Add(clientIP, limiter) // Add the new limiter to the cache
	return limiter
}

func limitMiddleware(next http.Handler, rateLimiter *RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr) // Extract client IP

		// Get the rate limiter for the client IP
		limiter := rateLimiter.GetLimiter(clientIP)

		// Check if the request is allowed
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests) // Respond with 429 Too Many Requests
			w.Write([]byte("Too Many Requests\n"))    // Optional: Add a message body
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

// logMiddleware is a middleware function that wraps an http.Handler.
// It logs details about each request, including the client IP, HTTP method, URL, status code, and response time.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()      // Record the start time of the request
		clientIP := r.RemoteAddr // Get the client's IP address
		method := r.Method       // Get the HTTP method (e.g., GET, POST)
		url := r.URL.String()    // Get the requested URL

		// Use a custom ResponseWriter to capture the status code
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r) // Pass the request to the next handler in the chain

		// Log the request details
		log.Printf("[%s] %s %s %s %d %s", start.Format(time.RFC3339), clientIP, method, url, lrw.statusCode, time.Since(start))
	})
}

// loggingResponseWriter is a custom implementation of http.ResponseWriter.
// It captures the HTTP status code for logging purposes.
type loggingResponseWriter struct {
	http.ResponseWriter     // Embeds the original ResponseWriter
	statusCode          int // Stores the HTTP status code
}

// WriteHeader overrides the WriteHeader method of http.ResponseWriter.
// It captures the status code and then calls the original WriteHeader method.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code                // Store the status code
	lrw.ResponseWriter.WriteHeader(code) // Call the original WriteHeader
}

// ServeProxy sets up and starts the reverse proxy server.
// It forwards requests to the specified target and logs each request using the middleware.
func ServeProxy(config []Config) {
	target := config[0].Target // Get the target URL from the config
	proxyPort := config[0].ProxyPort

	rateLimiter := NewRateLimiter(config[0]) // Create a new rate limiter instance

	// Parse the target URL
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing target host: %v", err)
	}

	// Create a reverse proxy for the target URL
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Wrap the proxy handler with the logging middleware
	handler := limitMiddleware(logMiddleware(proxy), rateLimiter)

	// Register the handler for the root path
	http.Handle("/", handler)

	// Log the server startup details
	log.Printf("Proxy server running on http://localhost:%d, forwarding to %s\n", proxyPort, target)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", proxyPort), nil))
}
