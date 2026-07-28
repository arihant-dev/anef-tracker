package auth

import (
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type RefreshPipeline = session.RefreshPipeline

func NewRefreshPipeline(client interface{}) *session.RefreshPipeline {
	return session.NewRefreshPipeline(nil)
}
