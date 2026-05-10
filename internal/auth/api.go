package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(header http.Header) (string, error) {

	fullAuthText := header.Get("Authorization")

	if fullAuthText == "" {
		return "", fmt.Errorf("no ApiKey token found in header")
	}

	fullAuthText = strings.TrimSpace(fullAuthText)
	textArr := strings.Fields(fullAuthText)
	var token string
	if len(textArr) <= 1 {
		return "", fmt.Errorf("not enough fields")
	}
	if textArr[0] != "ApiKey" {

		return "", fmt.Errorf("no apikey identifier found")
	}
	token = textArr[1]

	// trim away potential whitespace
	token = strings.TrimSpace(token)
	return token, nil
}
