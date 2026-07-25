package client_test

import (
	"context"
	"github.com/arihant-dev/anef-tracker/pkg/client"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClientDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	}))
	defer ts.Close()

	httpClient := client.NewHTTPClient(nil)
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("failed creating request: %v", err)
	}

	resp, err := httpClient.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("httpClient.Do failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", resp.StatusCode)
	}
}
