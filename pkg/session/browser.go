package session

import (
	"fmt"
	"net/http"
	"time"
)

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

func (b *BrowserAuthStrategy) Authenticate() (*Session, error) {
	if b.AuthToken == "" && len(b.Cookies) == 0 {
		return nil, fmt.Errorf("no browser cookies or tokens provided")
	}

	session := &Session{
		Provider:     "anef",
		URL:          "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour",
		Method:       "GET",
		Headers:      make(http.Header),
		AccessToken:  b.AuthToken,
		RefreshToken: b.RefreshToken,
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(1 * time.Hour),
		ImportSource: ImportBrowserAssisted,
	}

	for k, v := range b.Cookies {
		session.Cookies = append(session.Cookies, Cookie{Name: k, Value: v})
	}

	if b.AuthToken != "" {
		claims, err := DecodeJWT(b.AuthToken)
		if err == nil && claims != nil {
			session.User = claims.Login
			if claims.Exp > 0 {
				session.ExpiresAt = time.Unix(claims.Exp, 0)
			}
		}
	}

	return session, nil
}

func (b *BrowserAuthStrategy) Validate(s *Session) bool {
	return s != nil && s.AccessToken != "" && !s.IsExpired()
}
