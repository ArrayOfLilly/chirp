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
// fileserverHits: an atomic integer to safely track the number of hits to the file server in a concurrent environment.
// db: a pointer to a database.Queries object, which likely provides methods for interacting with the database.
// platform: a string representing the platform the API is running on.
// jwtSecret: a string containing the secret key used for signing JSON Web Tokens (JWTs).
// polkaKey: a string containing the Polka key (purpose not specified in this context).
type apiConfig struct {
	// safely incrementable int type for case of concurrent use
	fileserverHits 	atomic.Int32
	db 				*database.Queries
	platform       	string
	jwtSecret		string
	polkaKey		string
}

func main() {
	filepathRoot := "."

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

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable must be set")
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:		jwtSecret,
		polkaKey:		polkaKey,
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
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	// type HandlerFunc func(ResponseWriter, *Request)
	// The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers. 
	// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
	// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
	// ServeHTTP calls f(w, r).
	
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReady)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUserpdate)


	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpById)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUserUpgrade)
	
	// A Server defines parameters for running an HTTP server. 
	// The zero value for Server is a valid configuration.
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	// ListenAndServe listens on the TCP network address srv.Addr and 
	// then calls Serve to handle requests on incoming connections. 
	// opens a TCP socket
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}


