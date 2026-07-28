package auth

import (
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

func GetSessionFilePath() (string, error) {
	return session.GetSessionFilePath()
}
