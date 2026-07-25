package diff_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/diff"
	"testing"
)

func FuzzCompareSnapshots(f *testing.F) {
	f.Add([]byte(`{"statut":"INSTRUCTION_EN_COURS"}`), []byte(`{"statut":"TITRE_A_FABRIQUER"}`))
	f.Add([]byte(`{}`), []byte(`{}`))

	f.Fuzz(func(t *testing.T, bytesA, bytesB []byte) {
		// Verify diff engine does not panic on arbitrary snapshot fuzz inputs
		res, err := diff.CompareSnapshots(bytesA, bytesB)
		if err == nil && res != nil {
			_ = res.HasChanges
		}
	})
}
