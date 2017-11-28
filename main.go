package main

import (
	"log"
	"os"

	"git.starchasers.ovh/Starchasers/go-panel-microservice/microservice"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	service := microservice.NewPanelMicroservice(":"+port, "c2VjcmV0", "", "")
	service.Start()
}
