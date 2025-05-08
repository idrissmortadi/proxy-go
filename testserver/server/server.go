package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[TEST SERVER] %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	// Log headers
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}

	// Log body (if present)
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	if len(body) > 0 {
		fmt.Printf("  Body: %s\n", string(body))
	}

	w.Header().Set("X-Test-Server", "true")
	fmt.Fprint(w, "Hello from test server!\n")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Fprint(w, "Delayed response from test server")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Test server error", http.StatusInternalServerError)
}

func StartServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/slow", slowHandler)
	http.HandleFunc("/error", errorHandler)

	addr := ":3000"
	log.Printf("Test server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Test server failed: %v", err)
	}
}
