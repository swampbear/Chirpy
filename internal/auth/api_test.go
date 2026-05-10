package auth

import (
	"net/http"
	"testing"
)

func TestGetApiKey(t *testing.T) {

	// arrange
	header := http.Header{}
	header.Add("Authorization", "ApiKey  testToken     ")

	// act
	token, err := GetApiKey(header)
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
func TestMissingApiKey(t *testing.T) {

	//arrange
	header := http.Header{}

	// act
	_, err := GetApiKey(header)

	// assert
	if err == nil {
		t.Errorf("should not find ApiKey, but found: %s", err)
		return
	}
}
func TestWrongFormatApiKey(t *testing.T) {

	//arrange
	header := http.Header{}
	header.Add("Authorization", "ApiKey")

	// act
	_, err := GetApiKey(header)

	// assert
	if err == nil {
		t.Errorf("should not find ApiKey, but found: %s", err)
		return
	}
}
