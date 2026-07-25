package tui_test

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/tui/components"
	"testing"
	"time"
)

func TestPerformance10kItems(t *testing.T) {
	var items []string
	for i := 1; i <= 10000; i++ {
		items = append(items, fmt.Sprintf("Log Entry #%d: GET /api/sejour/usager/demande_sejour HTTP 200 (150ms)", i))
	}

	ds := &components.SliceDataSource{Items: items}
	vl := components.NewVirtualList(ds, 25)

	startTime := time.Now()
	_ = vl.RenderString()
	elapsed := time.Since(startTime)

	t.Logf("Rendered 10,000 rows viewport in %v", elapsed)

	if elapsed > 100*time.Millisecond {
		t.Errorf("rendering 10k items took longer than 100ms: %v", elapsed)
	}
}
