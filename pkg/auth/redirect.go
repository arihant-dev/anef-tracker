package auth

import (
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

func OpenBrowser(targetURL string) error {
	return session.OpenBrowser(targetURL)
}
