package auth

import (
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type BrowserAuthStrategy = session.BrowserAuthStrategy

func NewBrowserAuthStrategy(authToken, refreshToken string, cookies map[string]string) *BrowserAuthStrategy {
	return session.NewBrowserAuthStrategy(authToken, refreshToken, cookies)
}
