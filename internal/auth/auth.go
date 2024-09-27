package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const TokenTypeAccess TokenType = "chirpy-access"

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
var ErrBadAuthHeader = errors.New("malformed auth header")

// HashPassword generates a bcrypt hash from a given password.
//
// password is the string to be hashed.
// Returns the hashed password as a string and an error if the hashing process fails.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) 
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash checks if a given password matches a bcrypt hash.
//
// password is the string to be checked, and hash is the bcrypt hash to compare against.
// Returns an error if the password does not match the hash, or nil if they match.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT generates a JSON Web Token (JWT) for a given user ID, token secret, and expiration duration.
//
// userID is the unique identifier of the user, tokenSecret is the secret key used for signing the token, and expiresIn is the duration after which the token expires.
// Returns the signed JWT token as a string and an error if the signing process fails.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: 	string(TokenTypeAccess),
		IssuedAt: 	jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: 	jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: 	userID.String(),
	})
	return token.SignedString(signingKey)
}

// ValidateJWT validates a JSON Web Token (JWT) and returns the user ID.
//
// tokenString is the JWT token to be validated, and tokenSecret is the secret key used for signing the token.
// Returns the user ID as a UUID and an error if the validation fails.
// It is also takes account for the expiration date.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString, 
		&claimsStruct,
		func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	issuer, err := token.Claims.GetIssuer()
    if err != nil {
        return uuid.Nil, err
    }

    if issuer != string(TokenTypeAccess) {
        return uuid.Nil, errors.New("invalid issuer")
    }

	subject, err := token.Claims.GetSubject()
	if err != nil {
        return uuid.Nil, err
    }

	id, err := uuid.Parse(subject)
	if err != nil {
        return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
}


// GetBearerToken retrieves the bearer token from the Authorization header.
//
// headers: The HTTP headers.
// Returns the bearer token and an error if the header is missing or malformed.
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", ErrBadAuthHeader
	}
	return splitAuth[1], nil
}

// MakeRefreshToken generates a cryptographically secure random refresh token.
//
// No parameters.
// Returns the refresh token as a string and an error if the token generation fails.
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// GetApiKey retrieves the API key from the Authorization header.
//
// headers: The HTTP headers.
// Returns the API key and an error if the header is missing or malformed.
func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", ErrBadAuthHeader
	}
	return splitAuth[1], nil
}