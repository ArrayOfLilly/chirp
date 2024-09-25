package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// func handlerChirpCreate(w http.ResponseWriter, r *http.Request) {

// }

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	fmt.Println(params.Body)
	params.Body = filterProphane(params.Body, []string{
		"kerfuffle", "sharbert", "fornax"})
	fmt.Println(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: params.Body,
	})

}

func filterProphane(s string, wordsToFilter []string) string {
	words := strings.Split(s, " ")

	for i := range words {
		if contains(wordsToFilter, strings.ToLower(words[i])) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
			fmt.Println("match")
            return true
        }
    }
    return false
}