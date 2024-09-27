package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/ArrayOfLilly/chirp/internal/auth"
	"github.com/ArrayOfLilly/chirp/internal/database"
	"github.com/google/uuid"
)

// handlerChirpsCreate handles the creation of a new chirp.
//
// It takes an http.ResponseWriter and an http.Request as parameters.
// It returns no value, but writes the result of the creation (a chirp) to the http.ResponseWriter.
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

// validateChirp validates the length of a chirp message.
//
// It takes a string message as a parameter.
// Returns the validated message and an error if the message exceeds the maximum allowed length.
func validateChirp(msg string) (string, error) {
	const maxChirpLength = 140

	if len(msg) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	return msg, nil
}

// filterProphane filters out profane words from a given string.
//
// It takes a string `s` and a slice of words to filter `wordsToFilter` as parameters.
// It returns the filtered string with profane words replaced with "****".
func filterProphane(s string, wordsToFilter []string) string {
	words := strings.Split(s, " ")

	for i := range words {
		if contains(wordsToFilter, strings.ToLower(words[i])) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

// contains checks if a string slice contains a specific string element.
//
// It takes a string slice s and a string e as parameters.
// Returns a boolean indicating whether the string slice contains the string element.
func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// handlerGetAllChirps handles the retrieval of all chirps.
//
// It takes an http.ResponseWriter and an http.Request as parameters.
// Returns a JSON response containing a list of Chirp objects.
// If there's an optional author_id search query it lists the chirps by the userID of author
// otherwise it lists all chirps by creation time in ascending orded
// There's an optional sort querz string with asc and desc value for sorting
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	response := []Chirp{}

	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	for _, dbChirp := range dbChirps {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		response = append(response, databaseChirpToChirp(dbChirp))
	}

	sort.Slice(response, func(i, j int) bool {
		if sortDirection == "desc" {
			return response[i].CreatedAt.After(response[j].CreatedAt)
		}
		return response[i].CreatedAt.Before(response[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, response)
}
	

// handlerGetChirpById handles the retrieval of a chirp by its ID.
//
// It takes an http.ResponseWriter and an http.Request as parameters.
// Returns a JSON response containing a Chirp object by the pattern in the path.
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

	respondWithJSON(w, http.StatusOK, response{
		Chirp: databaseChirpToChirp(chirp),
	})
}

// handlerDeleteChirpById handles the deletion of a chirp by its ID.
//
// It takes an http.ResponseWriter and an http.Request as parameters.
// Returns no value, but writes the result of the deletion to the http.ResponseWriter.
func (cfg *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {

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

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Deletion is forbidden", fmt.Errorf("User is not the author of the chirp"))
		return
	}
	
	err = cfg.db.DeleteChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}