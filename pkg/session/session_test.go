package session

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSessionSaveAndLoad(t *testing.T) {
	sess := &Session{
		Provider:     "anef",
		URL:          "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour",
		Method:       "GET",
		AccessToken:  "test_access_token_123",
		RefreshToken: "test_refresh_token_456",
		User:         "9999999999",
		ImportSource: ImportCurl,
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(1 * time.Hour),
		Cookies: []Cookie{
			{Name: "auth_token", Value: "test_access_token_123"},
		},
		Headers: http.Header{
			"User-Agent": []string{"test-agent"},
		},
	}

	if err := SaveSession(sess); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	loaded, err := LoadSession()
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	if loaded.User != "9999999999" {
		t.Errorf("expected user 9999999999, got %s", loaded.User)
	}

	if loaded.AccessToken != "test_access_token_123" {
		t.Errorf("expected access token test_access_token_123, got %s", loaded.AccessToken)
	}

	if loaded.ImportSource != ImportCurl {
		t.Errorf("expected import source %s, got %s", ImportCurl, loaded.ImportSource)
	}
}

func TestPromptCurlPaste(t *testing.T) {
	input := "curl 'https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour' -H 'Authorization: Bearer token123'\n\n"
	buf := bytes.NewBufferString(input)

	sess, err := PromptCurlPaste(buf)
	if err != nil {
		t.Fatalf("PromptCurlPaste failed: %v", err)
	}

	if sess.AccessToken != "token123" {
		t.Errorf("expected access token token123, got %s", sess.AccessToken)
	}

	if sess.ImportSource != ImportCurl {
		t.Errorf("expected import source %s, got %s", ImportCurl, sess.ImportSource)
	}
}

func TestJWTDecoder(t *testing.T) {
	// Sample JWT header/payload/signature: {"login":"8888888888","exp":2000000000,"iss":"keycloak"}
	sampleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Ijg4ODg4ODg4ODgiLCJleHAiOjIwMDAwMDAwMDAsImlzcyI6ImtleWNsb2FrIn0.signature"

	claims, err := DecodeJWT(sampleToken)
	if err != nil {
		t.Fatalf("DecodeJWT failed: %v", err)
	}

	if claims.Login != "8888888888" {
		t.Errorf("expected login 8888888888, got %s", claims.Login)
	}

	if claims.Iss != "keycloak" {
		t.Errorf("expected issuer keycloak, got %s", claims.Iss)
	}
}

func TestInjectAuthHeaders(t *testing.T) {
	sess := &Session{
		AccessToken: "my_token",
		Cookies: []Cookie{
			{Name: "auth_token", Value: "my_token"},
		},
	}

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	InjectAuthHeaders(req, sess)

	if req.Header.Get("Authorization") != "Bearer my_token" {
		t.Errorf("unexpected authorization header: %s", req.Header.Get("Authorization"))
	}

	if !strings.Contains(req.Header.Get("Cookie"), "auth_token=my_token") {
		t.Errorf("cookie missing in request: %s", req.Header.Get("Cookie"))
	}
}

func TestValidateSession(t *testing.T) {
	sess := &Session{
		AccessToken: "valid_token",
		ExpiresAt:   time.Now().Add(2 * time.Hour),
		User:        "9999999999",
	}

	res := ValidateSession(sess)
	if !res.Ready {
		t.Errorf("expected session to be ready")
	}

	expiredSess := &Session{
		AccessToken: "expired_token",
		ExpiresAt:   time.Now().Add(-2 * time.Hour),
	}
	resExp := ValidateSession(expiredSess)
	if resExp.Ready {
		t.Errorf("expected expired session to not be ready")
	}
}

func TestRunSessionDoctor(t *testing.T) {
	rep := RunSessionDoctor()
	if len(rep.Checks) == 0 {
		t.Fatalf("expected doctor checks to run")
	}
}
