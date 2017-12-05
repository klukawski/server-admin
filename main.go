package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/microservice"
)

func handleTest(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Hello World!")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	service := microservice.NewPanelMicroservice(":"+port, "c2VjcmV0", "", "")
	service.Endpoints["/test"] = handleTest
	service.Start()
}
