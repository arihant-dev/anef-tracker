package auth

import (
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type JWTClaims = session.JWTClaims

func IsSessionExpired(s *session.Session) bool {
	return s.IsExpired()
}
