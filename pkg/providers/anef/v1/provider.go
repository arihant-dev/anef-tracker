package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/recorder"
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type ProviderV1 struct {
	HTTPClient *http.Client
	Recorder   *recorder.HTTPRecorder
	Session    *session.Session
}

func NewProviderV1(client *http.Client, rec *recorder.HTTPRecorder, sess *session.Session) *ProviderV1 {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &ProviderV1{
		HTTPClient: client,
		Recorder:   rec,
		Session:    sess,
	}
}

func (p *ProviderV1) Name() string {
	return "ANEF v1 API Provider"
}

func (p *ProviderV1) BaseURL() string {
	return "https://administration-etrangers-en-france.interieur.gouv.fr"
}

func (p *ProviderV1) Fetch(ctx context.Context) (*domain.Application, error) {
	// Standard residence permit dossier application endpoint
	targetURL := p.BaseURL() + "/api/sejour/usager/demande_sejour"

	// Only use session URL if it specifically targets a dossier or sejour resource
	if p.Session != nil && p.Session.URL != "" && strings.Contains(p.Session.URL, "/api/sejour/") {
		targetURL = p.Session.URL
	}

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating HTTP request: %w", err)
	}

	if p.Session != nil {
		session.InjectAuthHeaders(req, p.Session)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/27.0 Safari/605.1.15")
	}

	startTime := time.Now()
	resp, err := p.HTTPClient.Do(req)
	latency := time.Since(startTime).Milliseconds()

	if err != nil {
		if p.Recorder != nil {
			p.Recorder.Record(req, nil, nil, latency, err)
		}
		return nil, fmt.Errorf("ANEF HTTP request execution failed: %w", err)
	}

	var rawJSON []byte
	if p.Recorder != nil {
		rawJSON, err = p.Recorder.RecordAndRead(req, resp, latency)
	} else {
		defer resp.Body.Close()
		rawJSON, err = p.Recorder.RecordAndRead(req, resp, latency)
	}

	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ANEF API returned HTTP %d: %s", resp.StatusCode, string(rawJSON))
	}

	userLogin := ""
	if p.Session != nil {
		userLogin = p.Session.User
	}

	app, err := domain.MapJSONToApplication(rawJSON, userLogin)
	if err != nil {
		return nil, fmt.Errorf("failed mapping JSON to domain Application: %w", err)
	}

	return app, nil
}
