package anef

import (
	"net/http"

	v1 "github.com/arihant-dev/anef-tracker/pkg/providers/anef/v1"
	"github.com/arihant-dev/anef-tracker/pkg/recorder"
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type ANEFProvider = v1.ProviderV1

func NewANEFProvider(client *http.Client, rec *recorder.HTTPRecorder, sess *session.Session) *ANEFProvider {
	return v1.NewProviderV1(client, rec, sess)
}
