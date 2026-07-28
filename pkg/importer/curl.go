package importer

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var jwtPattern = regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)

type ImportedSession struct {
	Provider     string        `json:"provider"`
	URL          string        `json:"url,omitempty"`
	Method       string        `json:"method,omitempty"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Cookies      []Cookie      `json:"cookies"`
	Headers      http.Header   `json:"headers"`
	ExpiresAt    time.Time     `json:"expires_at"`
	IssuedAt     time.Time     `json:"issued_at"`
	User         string        `json:"user"`
	ImportSource ImportSource  `json:"import_source"`
	Capabilities Capabilities  `json:"capabilities"`
}

type CurlImporter struct{}

func NewCurlImporter() *CurlImporter {
	return &CurlImporter{}
}

func (c *CurlImporter) Import(raw []byte) (*ImportedSession, error) {
	curlCmd := string(raw)
	tokens := Tokenize(curlCmd)
	if len(tokens) == 0 {
		return nil, ErrEmptyInput
	}

	session := &ImportedSession{
		Provider:     "anef",
		Method:       "GET",
		Headers:      make(http.Header),
		Cookies:      make([]Cookie, 0),
		ImportSource: SourceCurl,
	}

	cookieMap := make(map[string]string)

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		if tok == "-H" || tok == "--header" {
			if i+1 < len(tokens) {
				i++
				headerVal := tokens[i]
				colonIdx := strings.Index(headerVal, ":")
				if colonIdx != -1 {
					k := strings.TrimSpace(headerVal[:colonIdx])
					v := strings.TrimSpace(headerVal[colonIdx+1:])
					session.Headers.Add(k, v)

					if strings.EqualFold(k, "Authorization") {
						session.AccessToken = extractTokenValue(v)
					}
					if strings.EqualFold(k, "Cookie") {
						ParseCookieHeader(v, cookieMap)
					}
				}
			}
			continue
		}

		if tok == "-b" || tok == "--cookie" {
			if i+1 < len(tokens) {
				i++
				ParseCookieHeader(tokens[i], cookieMap)
			}
			continue
		}

		if tok == "-X" || tok == "--request" {
			if i+1 < len(tokens) {
				i++
				session.Method = strings.ToUpper(tokens[i])
			}
			continue
		}

		if tok == "--url" {
			if i+1 < len(tokens) {
				i++
				session.URL = tokens[i]
			}
			continue
		}

		if strings.HasPrefix(tok, "http://") || strings.HasPrefix(tok, "https://") {
			if session.URL == "" {
				session.URL = tok
			}
		}
	}

	for k, v := range cookieMap {
		session.Cookies = append(session.Cookies, Cookie{Name: k, Value: v})
		if strings.EqualFold(k, "auth_token") && session.AccessToken == "" {
			session.AccessToken = v
		}
		if strings.EqualFold(k, "refresh_token") && session.RefreshToken == "" {
			session.RefreshToken = v
		}
	}

	if session.AccessToken == "" {
		if jwtMatch := jwtPattern.FindString(curlCmd); jwtMatch != "" {
			session.AccessToken = jwtMatch
		}
	}

	if session.AccessToken != "" {
		claims, err := DecodeJWT(session.AccessToken)
		if err == nil && claims != nil {
			session.User = claims.Login
			if claims.Exp > 0 {
				session.ExpiresAt = time.Unix(claims.Exp, 0)
			}
			if claims.Iat > 0 {
				session.IssuedAt = time.Unix(claims.Iat, 0)
			}
		}
	}

	if session.URL == "" {
		session.URL = "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour"
	}

	if session.AccessToken == "" && len(session.Cookies) == 0 {
		return nil, fmt.Errorf("no authentication token or cookies found in cURL command")
	}

	isExpired := !session.ExpiresAt.IsZero() && time.Now().After(session.ExpiresAt)
	session.Capabilities = ComputeCapabilities(session.AccessToken != "", session.RefreshToken != "", isExpired)

	return session, nil
}

func extractTokenValue(val string) string {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "Bearer ") {
		return strings.TrimPrefix(val, "Bearer ")
	}
	if strings.HasPrefix(val, "Token ") {
		return strings.TrimPrefix(val, "Token ")
	}
	return val
}
