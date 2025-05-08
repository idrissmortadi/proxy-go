package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ServeProxy(target string, port int) {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing target host: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received:", r.Method, r.URL)
		fmt.Println("Proxying to:", target)
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Starting proxy server on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
