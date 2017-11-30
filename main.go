package main

import (
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/microservice"
)

func handleRestart(w http.ResponseWriter, r *http.Request) {

}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	service := microservice.NewPanelMicroservice(":"+port, "c2VjcmV0", "", "")
	service.Start()
}
