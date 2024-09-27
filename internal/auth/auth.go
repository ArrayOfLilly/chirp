package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"strings"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string
const TokenTypeAccess TokenType = "chirpy-access"

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
var ErrBadAuthHeader = errors.New("malformed auth header")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) 
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: 	string(TokenTypeAccess),
		IssuedAt: 	jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: 	jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: 	fmt.Sprintf("%d", userID),
	})
	return token.SignedString(signingKey)
}

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
        return uuid.UUID{}, err
    }

    if issuer != string(TokenTypeAccess) {
        return uuid.UUID{}, errors.New("Invalid issuer")
    }

	subject, err := token.Claims.GetSubject()
	if err != nil {
        return uuid.UUID{}, err
    }

	subjectUUID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return subjectUUID, nil
}


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

