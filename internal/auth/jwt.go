package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// definte issued date and the expiration date
	currentTime := time.Now().UTC()
	currentDate := jwt.NewNumericDate(currentTime)
	expirationDate := jwt.NewNumericDate(currentTime.Add(expiresIn))
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  currentDate,
		ExpiresAt: expirationDate,
		Subject:   userID.String(),
	})

	secretBytes := []byte(tokenSecret)
	return token.SignedString(secretBytes)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	//Parse token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		secretBytes := []byte(tokenSecret)
		return secretBytes, nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parse token: %w", err)
	}
	userID, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("get subject: %w", err)
	}
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parse userID: %w", err)
	}
	return uuidUserID, nil
}
