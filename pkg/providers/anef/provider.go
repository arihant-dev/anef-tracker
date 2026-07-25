package anef

import (
	"github.com/arihant-dev/anef-tracker/pkg/auth"
	v1 "github.com/arihant-dev/anef-tracker/pkg/providers/anef/v1"
	"github.com/arihant-dev/anef-tracker/pkg/recorder"
	"net/http"
)

type ANEFProvider = v1.ProviderV1

func NewANEFProvider(client *http.Client, rec *recorder.HTTPRecorder, session *auth.CurlSession) *ANEFProvider {
	return v1.NewProviderV1(client, rec, session)
}
