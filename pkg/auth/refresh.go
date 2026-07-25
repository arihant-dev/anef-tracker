package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// RefreshPipeline handles automatic token renewal across strategies.
type RefreshPipeline struct {
	HTTPClient *http.Client
}

func NewRefreshPipeline(client *http.Client) *RefreshPipeline {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &RefreshPipeline{HTTPClient: client}
}

// RenewSession attempts to exchange the refresh_token for a new access_token.
func (p *RefreshPipeline) RenewSession(session *CurlSession) (*CurlSession, error) {
	if session == nil {
		return nil, fmt.Errorf("cannot renew nil session")
	}

	if session.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh_token present in session")
	}

	log.Println("[AUTH] Renewing OAuth2 access token via Keycloak endpoint...")
	tokenResp, err := RefreshToken(p.HTTPClient, session.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("token renewal failed: %w", err)
	}

	session.AuthToken = tokenResp.AccessToken
	session.RefreshToken = tokenResp.RefreshToken
	session.ExpiresAt = time.Now().Unix() + int64(tokenResp.ExpiresIn)
	session.Headers.Set("Authorization", "Token "+tokenResp.AccessToken)
	session.Cookies["auth_token"] = tokenResp.AccessToken
	session.Cookies["refresh_token"] = tokenResp.RefreshToken

	_ = SaveSession(session)
	log.Printf("[AUTH] Successfully renewed access token (expires in %d seconds)", tokenResp.ExpiresIn)

	return session, nil
}
