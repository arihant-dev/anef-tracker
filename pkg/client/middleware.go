package client

import (
	"github.com/arihant-dev/anef-tracker/pkg/auth"
	"github.com/arihant-dev/anef-tracker/pkg/log"
	"net/http"
)

// AuthMiddleware wraps an http.RoundTripper to inject session authentication and handle token refresh.
type AuthMiddleware struct {
	Next    http.RoundTripper
	Session *auth.CurlSession
}

func NewAuthMiddleware(next http.RoundTripper, session *auth.CurlSession) *AuthMiddleware {
	if next == nil {
		next = http.DefaultTransport
	}
	return &AuthMiddleware{
		Next:    next,
		Session: session,
	}
}

func (a *AuthMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	if a.Session != nil {
		// Check token expiration before sending request
		if a.Session.IsExpired() && a.Session.RefreshToken != "" {
			log.Info("Session token expired, attempting OAuth2 refresh", "user", a.Session.Login)
			pipeline := auth.NewRefreshPipeline(nil)
			if updatedSess, err := pipeline.RenewSession(a.Session); err == nil {
				a.Session = updatedSess
			} else {
				log.Warn("Session token renewal failed", "error", err)
			}
		}

		auth.InjectAuthHeaders(req, a.Session)
	}

	resp, err := a.Next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Retry on HTTP 401 Unauthorized if refresh_token present
	if resp.StatusCode == http.StatusUnauthorized && a.Session != nil && a.Session.RefreshToken != "" {
		log.Warn("Received HTTP 401, triggering auto-refresh retry...")
		resp.Body.Close()

		pipeline := auth.NewRefreshPipeline(nil)
		if updatedSess, err := pipeline.RenewSession(a.Session); err == nil {
			a.Session = updatedSess
			retryReq := req.Clone(req.Context())
			auth.InjectAuthHeaders(retryReq, a.Session)
			return a.Next.RoundTrip(retryReq)
		}
	}

	return resp, nil
}
