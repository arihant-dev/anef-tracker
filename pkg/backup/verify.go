package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func VerifyBackup(archivePath string) (*Manifest, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening backup file: %w", err)
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("invalid backup archive gzip format: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("invalid tar structure: %w", err)
		}

		if hdr.Name == "manifest.json" {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			var mf Manifest
			if err := json.Unmarshal(data, &mf); err != nil {
				return nil, fmt.Errorf("corrupted manifest json: %w", err)
			}
			return &mf, nil
		}
	}

	return nil, fmt.Errorf("manifest.json missing from backup archive")
}
