package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type RestoreResult struct {
	Manifest Manifest
	Restored bool
}

func RestoreBackup(archivePath string) (*RestoreResult, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening backup file '%s': %w", archivePath, err)
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed opening gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	res := &RestoreResult{}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading tar header: %w", err)
		}

		switch hdr.Name {
		case "manifest.json":
			data, _ := io.ReadAll(tr)
			_ = json.Unmarshal(data, &res.Manifest)
		case "anef.db":
			cwd, _ := os.Getwd()
			targetPath := filepath.Join(cwd, "data", "anef.db")
			_ = os.MkdirAll(filepath.Dir(targetPath), 0755)

			out, err := os.Create(targetPath)
			if err == nil {
				_, _ = io.Copy(out, tr)
				out.Close()
				res.Restored = true
			}
		}
	}

	return res, nil
}
