package main

import (
	"net/http"

	"./pkg/config"
	"./pkg/request"
	log "github.com/sirupsen/logrus"
)

const (
	configPath  = "./config.yml"
	secretsPath = "./secrets.yml"
)

func main() {
	// If this was a real application configurations should passed into our services via dependency injection
	config, err := config.NewConfig(configPath, secretsPath)
	if err != nil {
		log.Fatalf("couldn't load pkg path=%s, err=%s", configPath, err.Error())
	}
	requestHandler := request.NewRequestHandler(config)
	http.HandleFunc("/get-photos", requestHandler.HandleGetPhotos)
	http.HandleFunc("/get-all-photos", requestHandler.HandleGetAllPhotos)
	http.HandleFunc("/upload", requestHandler.HandleUpload)
	log.Info("Starting server on port 8090")

	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("Cannot start server err=%s", err)
	}
}
