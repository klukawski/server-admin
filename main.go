package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/klukawski/server-admin/auth"
)

var routes = flag.Bool("routes", false, "Generate router docs")

func main() {
	flag.Parse()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	if _, err := os.Stat("keys/local"); err != nil {
		auth.GenerateAndSaveKeypair()
	}

	auth.Config.Validator = jws.NewValidator(map[string]interface{}{"iss": "panel"}, time.Minute, time.Minute, nil)
	auth.Config.External = auth.LoadExternal("keys/external.pub")
	auth.Config.Local = auth.LoadPrivate("keys/local")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(auth.Auth)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})

	if *routes {
		fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/klukawski/server-admin",
			Intro:       "This is a schema of server-admin's REST API",
		}))
		return
	}

	http.ListenAndServe(":"+port, r)
}
