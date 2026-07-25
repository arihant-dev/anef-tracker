package auth

import (
	"fmt"
	"net/http"
	"time"
)

// BrowserAuthStrategy implements authentication via direct browser cookie & Bearer token injection.
type BrowserAuthStrategy struct {
	AuthToken    string
	RefreshToken string
	Cookies      map[string]string
}

func NewBrowserAuthStrategy(authToken, refreshToken string, cookies map[string]string) *BrowserAuthStrategy {
	if cookies == nil {
		cookies = make(map[string]string)
	}
	return &BrowserAuthStrategy{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		Cookies:      cookies,
	}
}

func (b *BrowserAuthStrategy) Name() string {
	return "Browser Strategy"
}

func (b *BrowserAuthStrategy) Authenticate() (*CurlSession, error) {
	if b.AuthToken == "" && len(b.Cookies) == 0 {
		return nil, fmt.Errorf("no browser cookies or tokens provided")
	}

	session := &CurlSession{
		URL:          "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour",
		Method:       "GET",
		Headers:      make(http.Header),
		Cookies:      b.Cookies,
		AuthToken:    b.AuthToken,
		RefreshToken: b.RefreshToken,
		IssuedAt:     time.Now().Unix(),
		ExpiresAt:    time.Now().Add(1 * time.Hour).Unix(),
	}

	if b.AuthToken != "" {
		claims, err := parseJWT(b.AuthToken)
		if err == nil {
			session.Login = claims.Login
			session.ExpiresAt = claims.Exp
		}
	}

	return session, nil
}

func (b *BrowserAuthStrategy) Validate(session *CurlSession) bool {
	return session != nil && session.AuthToken != "" && !session.IsExpired()
}

func (b *BrowserAuthStrategy) Refresh(session *CurlSession) (*CurlSession, error) {
	if session.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available in session")
	}
	tokenResp, err := RefreshToken(nil, session.RefreshToken)
	if err != nil {
		return nil, err
	}

	session.AuthToken = tokenResp.AccessToken
	session.RefreshToken = tokenResp.RefreshToken
	session.ExpiresAt = time.Now().Unix() + int64(tokenResp.ExpiresIn)
	session.Cookies["auth_token"] = tokenResp.AccessToken
	session.Cookies["refresh_token"] = tokenResp.RefreshToken

	return session, nil
}
