package main

import (
	"log"
	"net/http"
)

type Config struct{}

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: Router(),
	}

	log.Printf("Started server at addr: %v", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
