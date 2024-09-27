package main

import (
	"encoding/json"
	"net/http"

	"github.com/ArrayOfLilly/chirp/internal/auth"
	"github.com/ArrayOfLilly/chirp/internal/database"
)

// handlerUserCreate handles the user creation request.
//
// It expects a JSON payload in the request body with the fields "email" and "password".
// It returns a JSON response with the created user information.
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

// handlerUserpdate handles the user update request.
//
// It expects a JSON payload in the request body with the fields "email" and "password".
// It responds with a JSON payload containing the updated user information.
func (cfg *apiConfig) handlerUserpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 		string `json:"email"`
		Password 	string `json:"password"`
	}

	type response struct {
		User
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

	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user for token", err)
		return
	}

	var updatedEmail string
	var updatedPassword string

	if params.Email == "" {
		updatedEmail = user.Email
	} else {
		updatedEmail = params.Email
	}

	if params.Password == "" {
		updatedPassword = user.HashedPassword
	} else {
		updatedPassword, err = auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
			return
		}
	}

	updatedUser, err := cfg.db.UpdateUserData(r.Context(), database.UpdateUserDataParams{
		ID: 			user.ID,
		Email: 			updatedEmail,
		HashedPassword: updatedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: databaseUserToUser(updatedUser),
	})
}