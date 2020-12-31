package main

import (
	"./pkg/config"
	"log"
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
	/*limitValidator := limit_validator.NewLimitValidator(pkg)
	loadEvents, err := GetLoadEvents()
	if err != nil {
		log.Fatalf("couldn't load events err=%s", err)
	}
	validatedLoadEvents := limitValidator.ValidateLoadEvents(loadEvents)
	err = writeValidatedLoadEvents(validatedLoadEvents)
	if err != nil {
		log.Fatalf("couldn't write validated load events err=%s", err)
	}*/
}