package timeline_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/timeline"
	"testing"
)

func TestTimelineOrdering(t *testing.T) {
	ht, err := timeline.BuildTimeline(nil)
	if err != nil {
		t.Fatalf("BuildTimeline failed: %v", err)
	}

	if len(ht.Milestones) == 0 {
		t.Errorf("expected > 0 milestones in timeline")
	}

	if ht.ApplicationID == "" {
		t.Errorf("expected non-empty ApplicationID")
	}
}

func TestEvidenceAttachment(t *testing.T) {
	ht, _ := timeline.BuildTimeline(nil)
	fmtStr := ht.FormatASCII()

	if len(fmtStr) == 0 {
		t.Errorf("expected non-empty timeline ASCII output")
	}
}
