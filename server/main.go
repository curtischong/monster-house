package main

import (
	"./pkg/config"
	"log"
	"net/http"
	"./pkg/request"
)

const (
	configPath  = "./config.yaml"
)

func main() {
	// If this was a real application configurations should passed into our services via dependency injection
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load pkg path=%s, err=%s", configPath, err.Error())
	}
	print(config)
	requestHandler := request.NewRequestHandler()
	http.HandleFunc("/upload", requestHandler.HandleUpload)
	http.ListenAndServe(":8090", nil)
}