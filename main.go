package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// a struct that will hold any stateful, in-memory data we'll need to keep track of
type apiConfig struct {
	fileserverHits int
}

func main() {
	baseUrl := "."
	apiCfg := apiConfig{}

  	if err := godotenv.Load(); err != nil {
    	log.Fatal("Error loading .env file")
  	}

  	port := os.Getenv("PORT")

	// ServeMux is an HTTP request multiplexer. 
	// It matches the URL of each incoming request against a list of registered patterns and 
	// calls the handler for the pattern that most closely matches the URL.
	// The ServeMux is an http.Handler.
	mux := http.NewServeMux()

	// binds a handler against the route
	// basic handler,  which simply returns a built-in http.Handler
	
	// An http.Handler is just an interface: 
	// type Handler interface {
	// 	ServeHTTP(ResponseWriter, *Request)
	// }

	// FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.
	
	// To serve a directory on disk (/) under an alternate URL
	// path (/app), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(baseUrl)))))

	
	// http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))

	// type HandlerFunc func(ResponseWriter, *Request)
	// The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers. 
	// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
	// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
	// ServeHTTP calls f(w, r).
	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.HandleFunc("GET /api/error", handlerErr)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)


	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	
	
	// A Server defines parameters for running an HTTP server. 
	// The zero value for Server is a valid configuration.
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	// ListenAndServe listens on the TCP network address srv.Addr and 
	// then calls Serve to handle requests on incoming connections. 
	// opens a TCP socket
	log.Fatal(srv.ListenAndServe())
}