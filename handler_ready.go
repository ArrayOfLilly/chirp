package main

import (
	"net/http"
)

func handlerReady(w http.ResponseWriter, r *http.Request) {
	// add key, value to the header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// set the status code
	w.WriteHeader(http.StatusOK)
	// writes the response body in byte slice
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// func handlerErr(w http.ResponseWriter, r *http.Request) {
// 	// add key, value to the header
// 	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
// 	// set the status code
// 	w.WriteHeader(http.StatusInternalServerError)
// 	// writes the response body in byte slice
// 	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
// }