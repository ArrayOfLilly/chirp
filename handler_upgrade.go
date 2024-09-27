package main

import (
	"encoding/json"
	"net/http"

	"github.com/ArrayOfLilly/chirp/internal/auth"
	"github.com/google/uuid"
)

// handlerUserUpgrade handles the user upgrade request from the payment service provider, Polka.
// It is a webhhok.
// It returns no value, but writes the result of the upgrade to the http.ResponseWriter.
func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, r *http.Request) {

	type userData struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
  		Event string	`json:"event"`	
  		Data userData
  	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserById(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}