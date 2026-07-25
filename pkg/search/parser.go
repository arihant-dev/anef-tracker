package search

import (
	"strings"
	"time"
)

type ParsedQuery struct {
	RawText    string
	Field      string
	Type       string
	AfterDate  *time.Time
	SearchText string
}

func ParseQuery(queryStr string) *ParsedQuery {
	pq := &ParsedQuery{
		RawText: queryStr,
	}

	tokens := strings.Fields(queryStr)
	var textTokens []string

	for _, tok := range tokens {
		if strings.HasPrefix(tok, "field:") {
			pq.Field = strings.TrimPrefix(tok, "field:")
		} else if strings.HasPrefix(tok, "type:") {
			pq.Type = strings.TrimPrefix(tok, "type:")
		} else if strings.HasPrefix(tok, "after:") {
			dateStr := strings.TrimPrefix(tok, "after:")
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				pq.AfterDate = &t
			}
		} else {
			textTokens = append(textTokens, tok)
		}
	}

	pq.SearchText = strings.Join(textTokens, " ")
	return pq
}
