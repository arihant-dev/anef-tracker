package session

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type RefreshPipeline struct {
	HTTPClient *http.Client
}

func NewRefreshPipeline(client *http.Client) *RefreshPipeline {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &RefreshPipeline{HTTPClient: client}
}

func (p *RefreshPipeline) RenewSession(s *Session) (*Session, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot renew nil session")
	}

	if s.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh_token present in session")
	}

	log.Println("[SESSION] Renewing OAuth2 access token via Keycloak endpoint...")
	tokenResp, err := RefreshToken(p.HTTPClient, s.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("token renewal failed: %w", err)
	}

	s.AccessToken = tokenResp.AccessToken
	s.RefreshToken = tokenResp.RefreshToken
	s.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	s.Headers.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	// Update cookies
	foundAuth := false
	foundRef := false
	for i, c := range s.Cookies {
		if c.Name == "auth_token" {
			s.Cookies[i].Value = tokenResp.AccessToken
			foundAuth = true
		}
		if c.Name == "refresh_token" {
			s.Cookies[i].Value = tokenResp.RefreshToken
			foundRef = true
		}
	}
	if !foundAuth {
		s.Cookies = append(s.Cookies, Cookie{Name: "auth_token", Value: tokenResp.AccessToken})
	}
	if !foundRef && tokenResp.RefreshToken != "" {
		s.Cookies = append(s.Cookies, Cookie{Name: "refresh_token", Value: tokenResp.RefreshToken})
	}

	_ = SaveSession(s)
	log.Printf("[SESSION] Successfully renewed access token (expires in %d seconds)", tokenResp.ExpiresIn)

	return s, nil
}
