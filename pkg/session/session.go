package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var DefaultEncryptionKey = []byte("anef-tracker-secure-key-32bytes!")

const (
	ImportBrowserAssisted = "IMPORT_BROWSER_ASSISTED"
	ImportCurl            = "IMPORT_CURL"
	ImportPassword        = "IMPORT_PASSWORD"
	ImportFile            = "IMPORT_FILE"
)

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Capabilities struct {
	CanFetch        bool `json:"can_fetch"`
	CanDownload     bool `json:"can_download"`
	CanWatch        bool `json:"can_watch"`
	CanReplay       bool `json:"can_replay"`
	CanRefreshToken bool `json:"can_refresh_token"`
}

type Session struct {
	Provider     string            `json:"provider"`
	URL          string            `json:"url,omitempty"`
	Method       string            `json:"method,omitempty"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	Cookies      []Cookie          `json:"cookies"`
	Headers      http.Header       `json:"headers"`
	ExpiresAt    time.Time         `json:"expires_at"`
	IssuedAt     time.Time         `json:"issued_at"`
	User         string            `json:"user"`
	ImportSource string            `json:"import_source"`
	Capabilities Capabilities      `json:"capabilities"`

	// Legacy compatibility field mapping
	LegacyAuthToken string            `json:"auth_token,omitempty"`
	LegacyLogin     string            `json:"login,omitempty"`
	LegacyIat       int64             `json:"iat,omitempty"`
	LegacyExp       int64             `json:"exp,omitempty"`
	LegacyCookies   map[string]string `json:"legacy_cookies,omitempty"`
}

// UnmarshalJSON supports reading both old CurlSession format and new Session format.
func (s *Session) UnmarshalJSON(data []byte) error {
	type Alias Session
	aux := &struct {
		AuthToken string          `json:"auth_token"`
		Login     string          `json:"login"`
		Iat       int64           `json:"iat"`
		Exp       int64           `json:"exp"`
		RawCook   json.RawMessage `json:"cookies"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if s.AccessToken == "" && aux.AuthToken != "" {
		s.AccessToken = aux.AuthToken
	}
	if s.User == "" && aux.Login != "" {
		s.User = aux.Login
	}
	if s.ExpiresAt.IsZero() && aux.Exp > 0 {
		s.ExpiresAt = time.Unix(aux.Exp, 0)
	}
	if s.IssuedAt.IsZero() && aux.Iat > 0 {
		s.IssuedAt = time.Unix(aux.Iat, 0)
	}

	// Try unmarshaling cookies as map if slice is empty
	if len(s.Cookies) == 0 && len(aux.RawCook) > 0 {
		var cookMap map[string]string
		if err := json.Unmarshal(aux.RawCook, &cookMap); err == nil {
			for k, v := range cookMap {
				s.Cookies = append(s.Cookies, Cookie{Name: k, Value: v})
			}
		}
	}

	if s.Provider == "" {
		s.Provider = "anef"
	}
	if s.Headers == nil {
		s.Headers = make(http.Header)
	}

	return nil
}

// GetCookieMap returns a map representation of cookies.
func (s *Session) GetCookieMap() map[string]string {
	res := make(map[string]string)
	for _, c := range s.Cookies {
		res[c.Name] = c.Value
	}
	return res
}

// IsExpired checks if session token has expired.
func (s *Session) IsExpired() bool {
	if s == nil || s.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(s.ExpiresAt)
}

func GetSessionFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".anef")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "session.json"), nil
}

func SaveSession(s *Session) error {
	path, err := GetSessionFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	encrypted, err := encrypt(data, DefaultEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt session: %w", err)
	}

	return os.WriteFile(path, []byte(hex.EncodeToString(encrypted)), 0600)
}

func LoadSession() (*Session, error) {
	path, err := GetSessionFilePath()
	if err != nil {
		return nil, err
	}

	hexData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	encrypted, err := hex.DecodeString(string(hexData))
	if err != nil {
		return nil, fmt.Errorf("invalid hex session file: %w", err)
	}

	data, err := decrypt(encrypted, DefaultEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt session: %w", err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &sess, nil
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// InjectAuthHeaders adds current session tokens & cookies to an outgoing http.Request
func InjectAuthHeaders(req *http.Request, sess *Session) {
	if sess == nil {
		return
	}
	for k, vals := range sess.Headers {
		if strings.EqualFold(k, "Accept-Encoding") {
			continue
		}
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	if sess.AccessToken != "" {
		if !hasTokenPrefix(sess.AccessToken) {
			req.Header.Set("Authorization", "Bearer "+sess.AccessToken)
		} else {
			req.Header.Set("Authorization", sess.AccessToken)
		}
	}

	var cookieStrs []string
	for _, c := range sess.Cookies {
		cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	if len(cookieStrs) > 0 {
		req.Header.Set("Cookie", strings.Join(cookieStrs, "; "))
	}
}

func hasTokenPrefix(token string) bool {
	return len(token) > 6 && (token[:6] == "Token " || token[:7] == "Bearer ")
}
