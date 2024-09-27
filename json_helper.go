package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithError sends an error response to the client.
//
// It takes in the http.ResponseWriter, the HTTP status code, the error message, and the error.
// If the error is not nil, it logs the error.
// If the status code is greater than 499, it logs a message indicating a server error.
// It creates a JSON response with the error message and sends it to the client.
//
// Parameters:
//   - w: http.ResponseWriter to write the response.
//   - code: int representing the HTTP status code.
//   - msg: string representing the error message.
//   - err: error representing the error.
//
// Returns: None.
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

// respondWithJSON sends a JSON response to the client.
//
// It takes in the http.ResponseWriter to write the response, the HTTP status code, and the payload to be marshaled into JSON.
// Returns no value.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
