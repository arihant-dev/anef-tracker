package replay

import (
	"bytes"
	"fmt"
)

type MatchResult struct {
	Matched         bool   `json:"matched"`
	StatusCodeMatch bool   `json:"status_code_match"`
	HashMatch       bool   `json:"hash_match"`
	DiffSummary     string `json:"diff_summary"`
}

// CompareResponses checks status code and byte similarity between original and replayed responses.
func CompareResponses(origStatus, replayedStatus int, origBody, replayedBody []byte) *MatchResult {
	statusMatch := origStatus == replayedStatus
	hashMatch := bytes.Equal(origBody, replayedBody)
	matched := statusMatch && hashMatch

	summary := "Response payloads match exactly."
	if !matched {
		if !statusMatch {
			summary = fmt.Sprintf("Status code mismatch: HTTP %d vs HTTP %d.", origStatus, replayedStatus)
		} else {
			summary = fmt.Sprintf("Response body payload modified (Original: %d bytes, Replayed: %d bytes).", len(origBody), len(replayedBody))
		}
	}

	return &MatchResult{
		Matched:         matched,
		StatusCodeMatch: statusMatch,
		HashMatch:       hashMatch,
		DiffSummary:     summary,
	}
}
