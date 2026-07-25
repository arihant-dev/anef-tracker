package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuthMethod string

const (
	AuthMethodCurl     AuthMethod = "CURL"
	AuthMethodBrowser  AuthMethod = "BROWSER"
	AuthMethodKeycloak AuthMethod = "KEYCLOAK"
	AuthMethodBearer   AuthMethod = "BEARER"
)

// Strategy defines a pluggable authentication provider interface.
type Strategy interface {
	Name() string
	Authenticate() (*CurlSession, error)
	Validate(session *CurlSession) bool
	Refresh(session *CurlSession) (*CurlSession, error)
}

// AuthManager manages authentication strategies in cascade order.
type AuthManager struct {
	strategies []Strategy
}

func NewAuthManager(strategies ...Strategy) *AuthManager {
	return &AuthManager{strategies: strategies}
}

func (m *AuthManager) AuthenticateCascade() (*CurlSession, error) {
	var lastErr error
	for _, strat := range m.strategies {
		sess, err := strat.Authenticate()
		if err == nil && sess != nil {
			return sess, nil
		}
		lastErr = fmt.Errorf("[%s] %v", strat.Name(), err)
	}
	return nil, fmt.Errorf("all auth strategies failed: %w", lastErr)
}

// InjectAuthHeaders adds current session tokens & cookies to an outgoing http.Request
func InjectAuthHeaders(req *http.Request, sess *CurlSession) {
	if sess == nil {
		return
	}
	for k, vals := range sess.Headers {
		if strings.EqualFold(k, "Accept-Encoding") {
			continue // Omit explicit Accept-Encoding so net/http handles transparent decompression
		}
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	if sess.AuthToken != "" {
		if !hasTokenPrefix(sess.AuthToken) {
			req.Header.Set("Authorization", "Token "+sess.AuthToken)
		} else {
			req.Header.Set("Authorization", sess.AuthToken)
		}
	}

	var cookieStrs []string
	for k, v := range sess.Cookies {
		cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", k, v))
	}
	if len(cookieStrs) > 0 {
		req.Header.Set("Cookie", joinCookies(cookieStrs))
	}
}

func hasTokenPrefix(token string) bool {
	return len(token) > 6 && (token[:6] == "Token " || token[:7] == "Bearer ")
}

func joinCookies(cookies []string) string {
	res := ""
	for i, c := range cookies {
		if i > 0 {
			res += "; "
		}
		res += c
	}
	return res
}

func IsSessionExpired(sess *CurlSession) bool {
	if sess == nil || sess.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() >= sess.ExpiresAt
}
