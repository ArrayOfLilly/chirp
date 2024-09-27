package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArrayOfLilly/chirp/internal/auth"
	"github.com/ArrayOfLilly/chirp/internal/database"
)

// handlerLogin handles the login request.
//
// It expects a JSON payload in the request body with the fields "email" and "password".
// If the email and password are valid, it generates an access token and a refresh token.
// The access token is a JWT with a TTL of one hour.
// The refresh token is a random string.
// It responds with a JSON payload containing the user information, access token, and refresh token.
//
// Parameters:
//   - w: http.ResponseWriter to write the response.
//   - r: *http.Request containing the request.
//
// Returns: None.
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email 		string 	`json:"email"`
		Password 	string 	`json:"password"`
	}

	type response struct {
		User
		Token 			string `json:"token"`
		RefreshToken	string ` json:"refresh_token"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: 			databaseUserToUser(user),
		Token: 			accessToken,
		RefreshToken: 	refreshToken,
	})	
}