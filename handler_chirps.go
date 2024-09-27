package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ArrayOfLilly/chirp/internal/auth"
	"github.com/ArrayOfLilly/chirp/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	badWordList := []string{"kerfuffle", "sharbert", "fornax"}

	cleanedBody := filterProphane(params.Body, badWordList)
	validChirp, err := validateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
	}
	

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
        Body:   validChirp,
		UserID: userID,
    })
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
        return
    }

    respondWithJSON(w, http.StatusCreated, response{
		Chirp: databaseChirpToChirp(chirp),
	})
}

func validateChirp(msg string) (string, error) {
	const maxChirpLength = 140

	if len(msg) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	return msg, nil
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
            return true
        }
    }
    return false
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't succeed return chirps", err)
		return
	}

	response := []Chirp{}
	for _, chirp := range chirps {
		response = append(response, databaseChirpToChirp(chirp))
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {	
	
	type response struct {
		Chirp
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: databaseChirpToChirp(chirp),
	})
}