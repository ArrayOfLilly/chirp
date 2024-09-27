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
}

func databaseUserToUser(user database.User) User {
    return User{
        ID: 		user.ID,
        CreatedAt:	user.CreatedAt,
        UpdatedAt:	user.UpdatedAt,
        Email:		user.Email,
    }
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseChirpToChirp(chirp database.Chirp) Chirp {
	return Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}