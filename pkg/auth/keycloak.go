package auth

import (
	"net/http"

	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type KeycloakTokenResponse = session.KeycloakTokenResponse

func RefreshToken(httpClient *http.Client, refreshToken string) (*session.KeycloakTokenResponse, error) {
	return session.RefreshToken(httpClient, refreshToken)
}
