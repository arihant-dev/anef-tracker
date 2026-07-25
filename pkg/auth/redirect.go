package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// OpenBrowser opens the specified URL in the user's default web browser.
func OpenBrowser(targetURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{targetURL}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", targetURL}
	default:
		cmd = "xdg-open"
		args = []string{targetURL}
	}

	return exec.Command(cmd, args...).Start()
}

// StartLocalCallbackListener launches a temporary loopback HTTP server to capture authentication callbacks.
func StartLocalCallbackListener(port int, timeout time.Duration) (*CurlSession, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed binding callback port %d: %w", port, err)
	}

	var session *CurlSession
	var wg sync.WaitGroup
	wg.Add(1)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			curlParam := r.URL.Query().Get("curl")

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			if token != "" {
				session = &CurlSession{
					URL:       "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour",
					Method:    "GET",
					AuthToken: token,
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
					Headers:   make(http.Header),
					Cookies:   make(map[string]string),
				}
				session.Headers.Set("Authorization", "Bearer "+token)
				session.Cookies["auth_token"] = token
				if claims, err := parseJWT(token); err == nil && claims.Login != "" {
					session.Login = claims.Login
				}

				fmt.Fprint(w, `<!DOCTYPE html><html><head><title>ANEF Authentication</title><style>body{font-family:system-ui,sans-serif;background:#0f172a;color:#f8fafc;display:flex;justify-content:center;align-items:center;height:100vh;margin:0}.card{background:#1e293b;padding:2.5rem;border-radius:1rem;box-shadow:0 20px 25px -5px rgba(0,0,0,0.5);text-align:center;max-width:400px}h1{color:#38bdf8;margin-bottom:0.5rem}.badge{background:#0369a1;color:#e0f2fe;padding:0.25rem 0.75rem;border-radius:9999px;font-size:0.875rem}</style></head><body><div class="card"><h1>✓ Authentication Complete</h1><p>ANEF Tracker successfully captured your session.</p><p><span class="badge">You can close this window</span></p></div></body></html>`)
				wg.Done()
				return
			}

			if curlParam != "" {
				parsed, err := ParseCurl(curlParam)
				if err == nil {
					session = parsed
					fmt.Fprint(w, `<!DOCTYPE html><html><head><title>ANEF Authentication</title><style>body{font-family:system-ui,sans-serif;background:#0f172a;color:#f8fafc;display:flex;justify-content:center;align-items:center;height:100vh;margin:0}.card{background:#1e293b;padding:2.5rem;border-radius:1rem;box-shadow:0 20px 25px -5px rgba(0,0,0,0.5);text-align:center;max-width:400px}h1{color:#38bdf8;margin-bottom:0.5rem}.badge{background:#0369a1;color:#e0f2fe;padding:0.25rem 0.75rem;border-radius:9999px;font-size:0.875rem}</style></head><body><div class="card"><h1>✓ Session Imported</h1><p>ANEF Tracker parsed your cURL session successfully.</p><p><span class="badge">You can close this window</span></p></div></body></html>`)
					wg.Done()
					return
				}
			}

			// Helper bookmarklet page if user opens localhost callback directly
			fmt.Fprint(w, `<!DOCTYPE html><html><head><title>ANEF Tracker Callback</title><style>body{font-family:system-ui,sans-serif;background:#0f172a;color:#f8fafc;padding:2rem;line-height:1.6}code{background:#334155;padding:0.2rem 0.4rem;border-radius:0.25rem;color:#38bdf8}a.btn{display:inline-block;background:#0284c7;color:#fff;padding:0.6rem 1.2rem;border-radius:0.5rem;text-decoration:none;font-weight:bold;margin-top:1rem}</style></head><body><h2>ANEF Tracker Browser Authentication</h2><p>1. Open <a href="https://administration-etrangers-en-france.interieur.gouv.fr/usagers/" target="_blank" class="btn">ANEF Portal</a> and log in.</p><p>2. Paste token: <code>http://127.0.0.1:8484/callback?token=YOUR_JWT_TOKEN</code></p></body></html>`)
		}),
	}

	go func() {
		_ = server.Serve(listener)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		_ = server.Shutdown(ctx)
		return session, nil
	case <-ctx.Done():
		_ = server.Shutdown(context.Background())
		return nil, fmt.Errorf("browser authentication timed out after %s", timeout)
	}
}

// AuthenticateViaBrowser launches default browser and listens for session callback.
func AuthenticateViaBrowser(portalURL string, port int, timeout time.Duration) (*CurlSession, error) {
	if portalURL == "" {
		portalURL = "https://administration-etrangers-en-france.interieur.gouv.fr/usagers/"
	}
	if port <= 0 {
		port = 8484
	}
	if timeout <= 0 {
		timeout = 3 * time.Minute
	}

	fmt.Printf("Launching browser to %s ...\n", portalURL)
	if err := OpenBrowser(portalURL); err != nil {
		fmt.Printf("Notice: Could not automatically open browser (%v). Please open %s manually.\n", err, portalURL)
	}

	fmt.Printf("Listening for callback on http://127.0.0.1:%d/callback (Timeout: %s)...\n", port, timeout)
	return StartLocalCallbackListener(port, timeout)
}
