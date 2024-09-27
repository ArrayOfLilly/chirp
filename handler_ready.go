package main

import (
	"net/http"
)

// handlerReady handles the "/ready" endpoint and responds with a 200 OK status.
// It indicates the server up state.
func handlerReady(w http.ResponseWriter, r *http.Request) {
	// add key, value to the header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// set the status code
	w.WriteHeader(http.StatusOK)
	// writes the response body in byte slice
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
