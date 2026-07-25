package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CurlSession represents the extracted authentication session from a cURL command.
type CurlSession struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Headers      http.Header       `json:"headers"`
	Cookies      map[string]string `json:"cookies"`
	AuthToken    string            `json:"auth_token"`
	RefreshToken string            `json:"refresh_token"`
	Login        string            `json:"login"`
	IssuedAt     int64             `json:"iat"`
	ExpiresAt    int64             `json:"exp"`
}

// JWTClaims holds extracted JWT token fields.
type JWTClaims struct {
	Login     string `json:"login"`
	Iat       int64  `json:"iat"`
	Exp       int64  `json:"exp"`
	UrlPrefix string `json:"url_prefix"`
	Type      string `json:"type"`
	Fresh     bool   `json:"fresh"`
	Watcher   string `json:"watcher"`
}

// ParseCurl parses a cURL command string copied from Chrome DevTools.
func ParseCurl(curlCmd string) (*CurlSession, error) {
	session := &CurlSession{
		Headers: make(http.Header),
		Cookies: make(map[string]string),
		Method:  "GET",
	}

	lines := strings.Split(curlCmd, "\n")
	var fullCmd string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimSuffix(trimmed, "\\")
		trimmed = strings.TrimSpace(trimmed)
		fullCmd += " " + trimmed
	}

	// Extract URL
	urlStart := strings.Index(fullCmd, "curl '")
	if urlStart != -1 {
		urlStart += 6
		urlEnd := strings.Index(fullCmd[urlStart:], "'")
		if urlEnd != -1 {
			session.URL = fullCmd[urlStart : urlStart+urlEnd]
		}
	} else {
		// Fallback for curl "..."
		urlStart = strings.Index(fullCmd, `curl "`)
		if urlStart != -1 {
			urlStart += 6
			urlEnd := strings.Index(fullCmd[urlStart:], `"`)
			if urlEnd != -1 {
				session.URL = fullCmd[urlStart : urlStart+urlEnd]
			}
		}
	}

	// Extract headers
	headerParts := strings.Split(fullCmd, "-H '")
	if len(headerParts) == 1 {
		headerParts = strings.Split(fullCmd, `-H "`)
	}

	for _, part := range headerParts[1:] {
		endIdx := strings.Index(part, "'")
		if endIdx == -1 {
			endIdx = strings.Index(part, `"`)
		}
		if endIdx == -1 {
			continue
		}
		headerStr := part[:endIdx]
		colonIdx := strings.Index(headerStr, ":")
		if colonIdx != -1 {
			key := strings.TrimSpace(headerStr[:colonIdx])
			val := strings.TrimSpace(headerStr[colonIdx+1:])
			session.Headers.Add(key, val)

			if strings.EqualFold(key, "Authorization") {
				if strings.HasPrefix(val, "Token ") {
					session.AuthToken = strings.TrimPrefix(val, "Token ")
				} else if strings.HasPrefix(val, "Bearer ") {
					session.AuthToken = strings.TrimPrefix(val, "Bearer ")
				}
			}

			if strings.EqualFold(key, "Cookie") {
				parseCookies(val, session.Cookies)
			}
		}
	}

	// Check cookies for auth_token and refresh_token if not in headers
	if session.AuthToken == "" {
		if tok, ok := session.Cookies["auth_token"]; ok {
			session.AuthToken = tok
		} else if tok, ok := session.Cookies["Authorization"]; ok {
			session.AuthToken = tok
		}
	}

	if tok, ok := session.Cookies["refresh_token"]; ok {
		session.RefreshToken = tok
	}

	// Parse JWT claims from token
	if session.AuthToken != "" {
		claims, err := parseJWT(session.AuthToken)
		if err == nil {
			session.Login = claims.Login
			session.IssuedAt = claims.Iat
			session.ExpiresAt = claims.Exp
		}
	}

	if session.URL == "" {
		return nil, fmt.Errorf("failed to parse URL from cURL command")
	}

	return session, nil
}

func parseCookies(cookieStr string, target map[string]string) {
	parts := strings.Split(cookieStr, ";")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			val, err := url.QueryUnescape(kv[1])
			if err != nil {
				val = kv[1]
			}
			target[kv[0]] = val
		}
	}
}

func parseJWT(tokenStr string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid jwt format")
	}

	payload := parts[1]
	// Add base64 padding
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	data, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		data, err = base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, err
		}
	}

	var claims JWTClaims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// IsExpired checks if the JWT auth token has expired.
func (s *CurlSession) IsExpired() bool {
	if s.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() >= s.ExpiresAt
}
