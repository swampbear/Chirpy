package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	// arrange
	userID := uuid.New()

	tokenSecret := "verysecret"
	expiresIn, _ := time.ParseDuration("10m")
	//act
	tokenSignedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to get signed token: %s", err)
		return
	}
	validatedUserID, err := ValidateJWT(tokenSignedString, tokenSecret)

	//assert
	if err != nil {
		t.Errorf("failed to vaildate UUID %s", err)
		return
	}

	if validatedUserID != userID {
		t.Errorf("returned UUID does not match userId: %s", err)
		return
	}

}
func TestWrongSecret(t *testing.T) {
	// arrange
	userID := uuid.New()
	expiresIn := time.Minute * 10

	tokenWrongSecret := "verysecretbutwrong"
	tokenSecret := "verysecret"
	//act
	tokenSignedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to get signed token: %s", err)
		return
	}
	_, err = ValidateJWT(tokenSignedString, tokenWrongSecret)
	//assert
	if err == nil {
		t.Error("invalid token got past validation")
	}
}

func TestMalformedToken(t *testing.T) {
	_, err := ValidateJWT("really-illegal-string", "verysecret")
	if err == nil {
		t.Error("malformed token signed string got past validation")
	}
}

func TestExpiredToken(t *testing.T) {
	// arrange
	userID := uuid.New()
	tokenSecret := "verysecret"
	expiresIn := time.Second * -1

	//act
	tokenSignedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to get signed token: %s", err)
		return
	}
	_, err = ValidateJWT(tokenSignedString, tokenSecret)

	//assert
	if err == nil {
		t.Error("expired token got past validation")
	}
}
