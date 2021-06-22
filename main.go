package main

import (
	"log"
	"net/http"
	"time"

	"github.com/SimonTanner/hello-world-redis-app/api"
)

func main() {
	a := api.NewApi()

	server := &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
