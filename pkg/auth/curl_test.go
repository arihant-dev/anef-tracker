package auth_test

import (
	"testing"

	"github.com/arihant-dev/anef-tracker/pkg/session"
)

func TestParseCurl(t *testing.T) {
	curlSnippet := `curl 'https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour' \
-X 'GET' \
-H 'Pragma: no-cache' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Authorization: Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3ODQ5MjA1NzcsImV4cCI6MTc4NDkyMzI3NywibG9naW4iOiI5OTk5OTk5OTk5IiwidXJsX3ByZWZpeCI6Ii91c2FnZXIiLCJ0eXBlIjoiYXV0aCIsImZyZXNoIjp0cnVlLCJmcmVzaG5lc3NfZXhwIjoxNzg0OTIwNTc3LCJ3YXRjaGVyIjoiNDcxYzc0OWE0MzAyM2NmMGE5Yjk0ODI5MTExZThjYWJiZWM0MzYwOWM4NmRlYjhjYzEyYzJmYjUwOTg4OGFiMiIsImNsYWltcyI6eyJrY19mY19mbGFnIjpmYWxzZSwiZnVzaW9uX2ZsYWciOmZhbHNlfX0.xxx' \
-H 'Referer: https://administration-etrangers-en-france.interieur.gouv.fr/usagers/' \
-H 'Cookie: Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9; auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3ODQ5MjA1NzcsImV4cCI6MTc4NDkyMzI3NywibG9naW4iOiI5OTk5OTk5OTk5In0.xxx; refresh_token=mock_refresh_token; consentCookie=%7B%22necessary%22%3Atrue%7D'`

	sess, err := session.ParseCurl(curlSnippet)
	if err != nil {
		t.Fatalf("ParseCurl failed: %v", err)
	}

	if sess.URL != "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour" {
		t.Errorf("unexpected URL: %s", sess.URL)
	}

	if sess.User != "9999999999" {
		t.Errorf("expected user 9999999999, got %s", sess.User)
	}

	if sess.RefreshToken != "mock_refresh_token" {
		t.Errorf("expected refresh_token mock_refresh_token, got %s", sess.RefreshToken)
	}
}
