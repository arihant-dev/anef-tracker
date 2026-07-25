package auth_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/auth"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestLocalCallbackListener(t *testing.T) {
	port := 8499
	doneChan := make(chan *auth.CurlSession, 1)

	go func() {
		sess, err := auth.StartLocalCallbackListener(port, 3*time.Second)
		if err == nil {
			doneChan <- sess
		} else {
			doneChan <- nil
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Simulate callback request
	callbackURL := "http://127.0.0.1:8499/callback?token=TEST_JWT_TOKEN_12345"
	resp, err := http.Get(callbackURL)
	if err != nil {
		t.Fatalf("Callback request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", resp.StatusCode)
	}

	select {
	case sess := <-doneChan:
		if sess == nil {
			t.Fatalf("expected non-nil session")
		}
		if sess.AuthToken != "TEST_JWT_TOKEN_12345" {
			t.Errorf("expected AuthToken TEST_JWT_TOKEN_12345, got %s", sess.AuthToken)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("callback listener timed out")
	}
}

func TestLocalCallbackListenerCurlParam(t *testing.T) {
	port := 8498
	doneChan := make(chan *auth.CurlSession, 1)

	go func() {
		sess, err := auth.StartLocalCallbackListener(port, 3*time.Second)
		if err == nil {
			doneChan <- sess
		} else {
			doneChan <- nil
		}
	}()

	time.Sleep(100 * time.Millisecond)

	curlCmd := "curl 'https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour' -H 'Authorization: Bearer TEST_BEARER'"
	callbackURL := "http://127.0.0.1:8498/callback?curl=" + url.QueryEscape(curlCmd)

	resp, err := http.Get(callbackURL)
	if err != nil {
		t.Fatalf("Callback request failed: %v", err)
	}
	defer resp.Body.Close()

	select {
	case sess := <-doneChan:
		if sess == nil {
			t.Fatalf("expected non-nil session")
		}
		if sess.AuthToken != "TEST_BEARER" {
			t.Errorf("expected AuthToken TEST_BEARER, got %s", sess.AuthToken)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("callback listener timed out")
	}
}
