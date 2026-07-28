package session

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	KeycloakTokenEndpoint = "https://sso.anef.dgef.interieur.gouv.fr/auth/realms/anef-usagers/protocol/openid-connect/token"
	KeycloakClientID      = "anef-usagers"
)

type KeycloakTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func AuthenticateWithCredentials(httpClient *http.Client, username, password string) (*Session, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	form := url.Values{}
	form.Set("client_id", KeycloakClientID)
	form.Set("grant_type", "password")
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", "openid roles web-origins email basic profile")

	req, err := http.NewRequest("POST", KeycloakTokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/27.0 Safari/605.1.15")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth network request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp KeycloakTokenResponse
		_ = json.Unmarshal(body, &errResp)
		if errResp.Error == "unauthorized_client" || resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("direct password login disabled by ANEF portal (Keycloak 401 unauthorized_client).\n\n" +
				"ANEF security policies disable direct password token grants. Please import session via cURL:\n" +
				"  1. Open Chrome/Firefox DevTools (F12) -> Network tab\n" +
				"  2. Log into https://administration-etrangers-en-france.interieur.gouv.fr\n" +
				"  3. Right-click any API request -> Copy -> Copy as cURL\n" +
				"  4. Run: anef login --curl \"<paste_curl_here>\"")
		}
		if errResp.ErrorDescription != "" {
			return nil, fmt.Errorf("keycloak error (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.ErrorDescription)
		}
		return nil, fmt.Errorf("keycloak error HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp KeycloakTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	sess := &Session{
		URL:          "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour",
		Method:       "GET",
		Headers:      make(http.Header),
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		User:         username,
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		ImportSource: "credentials",
		Cookies:      []Cookie{{Name: "auth_token", Value: tokenResp.AccessToken}},
	}

	if tokenResp.RefreshToken != "" {
		sess.Cookies = append(sess.Cookies, Cookie{Name: "refresh_token", Value: tokenResp.RefreshToken})
	}

	sess.Headers.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	return sess, nil
}

func RefreshToken(httpClient *http.Client, refreshToken string) (*KeycloakTokenResponse, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	form := url.Values{}
	form.Set("client_id", KeycloakClientID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", KeycloakTokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/27.0 Safari/605.1.15")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh network request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp KeycloakTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed parsing token refresh response: %w", err)
	}

	return &tokenResp, nil
}
