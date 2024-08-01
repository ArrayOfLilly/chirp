package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// baseUrl = "."

	
  	if err := godotenv.Load(); err != nil {
    	log.Fatal("Error loading .env file")
  	}

  	port := os.Getenv("PORT")

	// ServeMux is an HTTP request multiplexer. 
	// It matches the URL of each incoming request against a list of registered patterns and 
	// calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()

	// A Server defines parameters for running an HTTP server. 
	// The zero value for Server is a valid configuration.
	srv := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	// ListenAndServe listens on the TCP network address srv.Addr and 
	// then calls Serve to handle requests on incoming connections. 
	log.Fatal(srv.ListenAndServe())
}