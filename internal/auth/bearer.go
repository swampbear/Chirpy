package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(header http.Header) (string, error) {
	fullAuthText := header.Get("Authorization")

	if fullAuthText == "" {
		return "", fmt.Errorf("no bearer token found in header")
	}

	fullAuthText = strings.TrimSpace(fullAuthText)
	textArr := strings.Fields(fullAuthText)
	var token string
	if len(textArr) <= 1 {
		return "", fmt.Errorf("recieved wrong format bearertoken")
	}
	token = textArr[1]

	// trim away potential whitespace
	token = strings.TrimSpace(token)
	return token, nil
}
