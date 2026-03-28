package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header")
	}

	bearer := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	return bearer, nil
}
