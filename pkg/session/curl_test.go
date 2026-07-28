package session

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "Simple single line",
			input: `curl 'https://api.example.com' -H 'Authorization: Bearer token123'`,
			expected: []string{
				"curl",
				"https://api.example.com",
				"-H",
				"Authorization: Bearer token123",
			},
		},
		{
			name: "Multiline bash with backslashes",
			input: `curl 'https://api.example.com' \
  -H 'Authorization: Bearer token123' \
  -H 'Cookie: auth_token=abc'`,
			expected: []string{
				"curl",
				"https://api.example.com",
				"-H",
				"Authorization: Bearer token123",
				"-H",
				"Cookie: auth_token=abc",
			},
		},
		{
			name: "PowerShell tick continuations",
			input: `curl "https://api.example.com" ` + "`" + `
  -H "Authorization: Bearer token123"`,
			expected: []string{
				"curl",
				"https://api.example.com",
				"-H",
				"Authorization: Bearer token123",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := Tokenize(tc.input)
			if len(tokens) != len(tc.expected) {
				t.Fatalf("expected %d tokens, got %d: %v", len(tc.expected), len(tokens), tokens)
			}
			for i := range tokens {
				if tokens[i] != tc.expected[i] {
					t.Errorf("token %d: expected %q, got %q", i, tc.expected[i], tokens[i])
				}
			}
		})
	}
}

func TestParseCurl(t *testing.T) {
	rawCurl := `curl 'https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Ijk5OTk5OTk5OTkiLCJpYXQiOjE2MDAwMDAwMDAsImV4cCI6MjAwMDAwMDAwMH0.signature' \
  -H 'Cookie: auth_token=abc; refresh_token=xyz'`

	sess, err := ParseCurl(rawCurl)
	if err != nil {
		t.Fatalf("ParseCurl failed: %v", err)
	}

	if sess.AccessToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Ijk5OTk5OTk5OTkiLCJpYXQiOjE2MDAwMDAwMDAsImV4cCI6MjAwMDAwMDAwMH0.signature" {
		t.Errorf("unexpected access token: %s", sess.AccessToken)
	}
	if sess.User != "9999999999" {
		t.Errorf("expected user 9999999999, got %s", sess.User)
	}
	if sess.ImportSource != ImportCurl {
		t.Errorf("expected import source %s, got %s", ImportCurl, sess.ImportSource)
	}
}
