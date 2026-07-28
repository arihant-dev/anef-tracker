package session

import (
	"time"

	"github.com/arihant-dev/anef-tracker/pkg/importer"
)

type JWTClaims = importer.JWTClaims

func DecodeJWT(tokenStr string) (*JWTClaims, error) {
	return importer.DecodeJWT(tokenStr)
}

func GetUserFromJWT(token string) string {
	return importer.GetUserFromJWT(token)
}

func GetExpiryFromJWT(token string) time.Time {
	return importer.GetExpiryFromJWT(token)
}
