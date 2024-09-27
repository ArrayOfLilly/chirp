package main

import (
	"time"

	"github.com/ArrayOfLilly/chirp/internal/database"
	"github.com/google/uuid"
)

type User struct {
    ID        		uuid.UUID	`json:"id"`
    CreatedAt 		time.Time	`json:"created_at"`
    UpdatedAt 		time.Time	`json:"updated_at"`
    Email      		string		`json:"email"`
	IsChirpyRed		bool		`json:"is_chirpy_red"`
}

// databaseUserToUser converts a database.User object to a User object.
//
// It takes a database.User object as a parameter.
// Returns a User object.
func databaseUserToUser(user database.User) User {
    return User{
        ID: 		 user.ID,
        CreatedAt:	 user.CreatedAt,
        UpdatedAt:	 user.UpdatedAt,
        Email:		 user.Email,
		IsChirpyRed: user.IsChirpyRed,
    }
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// databaseChirpToChirp converts a database.Chirp object to a Chirp object.
//
// It takes a database.Chirp object as a parameter.
// Returns a Chirp object.
func databaseChirpToChirp(chirp database.Chirp) Chirp {
	return Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}