package replay_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/replay"
	"testing"
)

func TestResponseMatching(t *testing.T) {
	origBody := []byte(`{"statut":"TITRE_A_FABRIQUER"}`)
	replayedBody := []byte(`{"statut":"TITRE_A_FABRIQUER"}`)

	match := replay.CompareResponses(200, 200, origBody, replayedBody)
	if !match.Matched {
		t.Errorf("expected responses to match")
	}

	diffBody := []byte(`{"statut":"TITRE_DISPONIBLE"}`)
	mismatch := replay.CompareResponses(200, 200, origBody, diffBody)
	if mismatch.Matched {
		t.Errorf("expected responses to mismatch")
	}
}

func TestHashBytes(t *testing.T) {
	h1 := replay.HashBytes([]byte("hello"))
	h2 := replay.HashBytes([]byte("hello"))
	if h1 != h2 {
		t.Errorf("hashes should be identical")
	}
}
