package search_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/search"
	"testing"
	"time"
)

func TestFieldSearch(t *testing.T) {
	pq := search.ParseQuery("field:adresse")
	if pq.Field != "adresse" {
		t.Errorf("expected Field 'adresse', got '%s'", pq.Field)
	}

	f1 := domain.FieldObservation{Path: "informations_personnelles.adresse.ville"}
	f2 := domain.FieldObservation{Path: "statut"}

	if !search.MatchFieldObservation(f1, pq) {
		t.Errorf("expected f1 to match field:adresse query")
	}

	if search.MatchFieldObservation(f2, pq) {
		t.Errorf("expected f2 not to match field:adresse query")
	}
}

func TestCombinedFilters(t *testing.T) {
	pq := search.ParseQuery("type:FIELD_DISCOVERED field:statut")
	if pq.Type != "FIELD_DISCOVERED" || pq.Field != "statut" {
		t.Errorf("unexpected parsed query fields")
	}

	ev1 := domain.Event{
		Type:      "FIELD_DISCOVERED",
		FieldPath: "statut",
		Timestamp: time.Now(),
	}

	ev2 := domain.Event{
		Type:      "STATUS_CHANGE",
		FieldPath: "statut",
		Timestamp: time.Now(),
	}

	if !search.MatchEvent(ev1, pq) {
		t.Errorf("expected ev1 to match combined query")
	}

	if search.MatchEvent(ev2, pq) {
		t.Errorf("expected ev2 to fail type filter")
	}
}
