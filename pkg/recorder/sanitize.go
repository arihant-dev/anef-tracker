package recorder

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// SanitizeHeaders returns a copy of http.Header with secrets redacted.
func SanitizeHeaders(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	clean := h.Clone()
	for k := range clean {
		lowerKey := strings.ToLower(k)
		if lowerKey == "authorization" || lowerKey == "cookie" || lowerKey == "set-cookie" || lowerKey == "x-auth-token" {
			clean.Set(k, "[REDACTED SECRET]")
		}
	}
	return clean
}

// SanitizeString returns a string with token/cookie substrings redacted.
func SanitizeString(input string) string {
	if strings.Contains(input, "Authorization:") {
		input = "[REDACTED AUTHORIZATION HEADER]"
	}
	return input
}

func readBodyAndRestore(body io.ReadCloser) ([]byte, io.ReadCloser, error) {
	if body == nil {
		return nil, nil, nil
	}
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}
	return data, io.NopCloser(bytes.NewBuffer(data)), nil
}
