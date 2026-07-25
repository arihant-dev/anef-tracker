package report

import (
	"encoding/json"
)

func (r *EvidenceReport) RenderJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
