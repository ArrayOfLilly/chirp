package main

import (
	"encoding/json"
	"net/http"

	"github.com/ArrayOfLilly/chirp/internal/auth"
	"github.com/ArrayOfLilly/chirp/internal/database"
)

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 		string `json:"email"`
		Password 	string `json:"password"`
	}

	type response struct {
		User
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: databaseUserToUser(user),
	})
}