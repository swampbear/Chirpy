package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(header http.Header) (string, error) {
	fullAuthText := header.Get("Authorization")

	if fullAuthText == "" {
		return "", fmt.Errorf("no bearer token found")
	}

	token := strings.Split(fullAuthText, " ")[1]
	// trim away potential whitespace
	token = strings.TrimSpace(token)
	return token, nil
}
