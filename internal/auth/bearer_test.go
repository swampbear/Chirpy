package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {

	// arrange
	header := http.Header{}
	header.Add("Authorization", "Bearer testToken")

	// act
	token, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// assert
	expectedToken := "testToken"
	if token != expectedToken {
		t.Errorf("should get: %s, got %s", expectedToken, token)
	}
}
func TestMissingBearerToken(t *testing.T) {

	//arrange
	header := http.Header{}

	// act
	_, err := GetBearerToken(header)

	// assert
	if err == nil {
		t.Errorf("should not find bearertoken, but found: %s", err)
		return
	}
}
