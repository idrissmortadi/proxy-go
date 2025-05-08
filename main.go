package main

import (
	"flag"

	"github.com/idrissmortadi/proxy-go/proxy"
)

func main() {
	target := flag.String("target", "http://localhost:3000", "Target server URL")
	port := flag.Int("port", 8080, "Proxy server port")
	flag.Parse()

	proxy.ServeProxy(*target, *port)
}
