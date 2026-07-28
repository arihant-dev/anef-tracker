package session

import (
	"github.com/arihant-dev/anef-tracker/pkg/importer"
)

// ParseCurl parses a raw cURL command string using the dedicated importer package.
func ParseCurl(curlCmd string) (*Session, error) {
	imp := importer.NewCurlImporter()
	imported, err := imp.Import([]byte(curlCmd))
	if err != nil {
		return nil, err
	}

	sess := &Session{
		Provider:     imported.Provider,
		URL:          imported.URL,
		Method:       imported.Method,
		AccessToken:  imported.AccessToken,
		RefreshToken: imported.RefreshToken,
		ExpiresAt:    imported.ExpiresAt,
		IssuedAt:     imported.IssuedAt,
		User:         imported.User,
		ImportSource: string(imported.ImportSource),
		Headers:      imported.Headers,
		Capabilities: Capabilities{
			CanFetch:        imported.Capabilities.CanFetch,
			CanDownload:     imported.Capabilities.CanDownload,
			CanWatch:        imported.Capabilities.CanWatch,
			CanReplay:       imported.Capabilities.CanReplay,
			CanRefreshToken: imported.Capabilities.CanRefreshToken,
		},
	}

	for _, c := range imported.Cookies {
		sess.Cookies = append(sess.Cookies, Cookie{Name: c.Name, Value: c.Value})
	}

	return sess, nil
}
