package main

import (
	"log"
	"os"

	"github.com/idrissmortadi/proxy-go/proxy"
	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.Open("config.yaml")
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
