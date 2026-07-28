package importer

import (
	"fmt"
)

// Importer defines the unified interface for extracting an authenticated session from raw bytes or strings.
type Importer interface {
	Import(raw []byte) (interface{}, error)
}

type ImportSource string

const (
	SourceBrowserAssisted ImportSource = "IMPORT_BROWSER_ASSISTED"
	SourceCurl            ImportSource = "IMPORT_CURL"
	SourcePassword        ImportSource = "IMPORT_PASSWORD"
	SourceFile            ImportSource = "IMPORT_FILE"
)

type Capabilities struct {
	CanFetch        bool `json:"can_fetch"`
	CanDownload     bool `json:"can_download"`
	CanWatch        bool `json:"can_watch"`
	CanReplay       bool `json:"can_replay"`
	CanRefreshToken bool `json:"can_refresh_token"`
}

func ComputeCapabilities(hasAccessToken, hasRefreshToken bool, isExpired bool) Capabilities {
	valid := hasAccessToken && !isExpired
	return Capabilities{
		CanFetch:        valid,
		CanDownload:     valid,
		CanWatch:        valid,
		CanReplay:       valid,
		CanRefreshToken: hasRefreshToken,
	}
}

var ErrEmptyInput = fmt.Errorf("empty import payload")
