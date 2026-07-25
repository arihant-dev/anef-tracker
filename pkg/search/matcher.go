package search

import (
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"strings"
)

func MatchEvent(ev domain.Event, pq *ParsedQuery) bool {
	if pq == nil {
		return true
	}

	if pq.Type != "" && !strings.EqualFold(ev.Type, pq.Type) {
		return false
	}

	if pq.Field != "" && !strings.Contains(strings.ToLower(ev.FieldPath), strings.ToLower(pq.Field)) {
		return false
	}

	if pq.AfterDate != nil && ev.Timestamp.Before(*pq.AfterDate) {
		return false
	}

	if pq.SearchText != "" {
		st := strings.ToLower(pq.SearchText)
		if !strings.Contains(strings.ToLower(ev.FieldPath), st) &&
			!strings.Contains(strings.ToLower(ev.Type), st) &&
			!strings.Contains(strings.ToLower(ev.OldVal), st) &&
			!strings.Contains(strings.ToLower(ev.NewVal), st) {
			return false
		}
	}

	return true
}

func MatchFieldObservation(f domain.FieldObservation, pq *ParsedQuery) bool {
	if pq == nil {
		return true
	}

	if pq.Field != "" && !strings.Contains(strings.ToLower(f.Path), strings.ToLower(pq.Field)) {
		return false
	}

	if pq.AfterDate != nil && f.FirstSeen.Before(*pq.AfterDate) {
		return false
	}

	if pq.SearchText != "" {
		st := strings.ToLower(pq.SearchText)
		if !strings.Contains(strings.ToLower(f.Path), st) &&
			!strings.Contains(strings.ToLower(f.Type), st) &&
			!strings.Contains(strings.ToLower(f.Endpoint), st) {
			return false
		}
	}

	return true
}
