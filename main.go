package main

import (
	"flag"
	"log"
	"os"

	"github.com/idrissmortadi/proxy-go/proxy"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config []proxy.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding YAML: %v", err)
	}

	proxy.ServeProxy(config)
}
