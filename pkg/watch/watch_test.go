package watch_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/watch"
	"testing"
)

func TestChangeDetection(t *testing.T) {
	err := watch.RecordWatchRun(nil, "SUCCESS", 0)
	if err != nil {
		t.Errorf("expected clean RecordWatchRun execution, got %v", err)
	}
}
