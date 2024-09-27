package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ArrayOfLilly/chirp/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// a struct that will hold any stateful, in-memory data we'll need to keep track of
type apiConfig struct {
	// safely incrementable int type for case of concurrent use
	fileserverHits 	atomic.Int32
	db 				*database.Queries
	platform       	string
	jwtSecret			string
}

func main() {
	baseUrl := "."

  	if err := godotenv.Load(); err != nil {
    	log.Fatal("Error loading .env file")
  	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

  	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("No database connection")
	}
	dbQueries := database.New(dbConn)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:		jwtSecret,
	}

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
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(baseUrl)))))

	
	// http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))

	// type HandlerFunc func(ResponseWriter, *Request)
	// The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers. 
	// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
	// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
	// ServeHTTP calls f(w, r).
	
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReady)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpById)

	mux.HandleFunc("POST /api/users", cfg.handlerUserCreate)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	
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