package components_test

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/tui/components"
	"testing"
)

func TestVirtualListScrolling(t *testing.T) {
	var items []string
	for i := 1; i <= 100; i++ {
		items = append(items, fmt.Sprintf("Item %d", i))
	}

	ds := &components.SliceDataSource{Items: items}
	vl := components.NewVirtualList(ds, 10)

	if vl.Offset != 0 {
		t.Errorf("expected initial offset 0, got %d", vl.Offset)
	}

	vl.ScrollDown(15)
	if vl.Offset != 15 {
		t.Errorf("expected offset 15 after ScrollDown(15), got %d", vl.Offset)
	}

	visible := vl.RenderVisible()
	if len(visible) != 10 {
		t.Errorf("expected 10 visible items, got %d", len(visible))
	}

	if visible[0] != "Item 16" {
		t.Errorf("expected first visible item 'Item 16', got '%s'", visible[0])
	}
}

func TestVirtualListRendering(t *testing.T) {
	var items []string
	for i := 1; i <= 5; i++ {
		items = append(items, fmt.Sprintf("Row %d", i))
	}

	ds := &components.SliceDataSource{Items: items}
	vl := components.NewVirtualList(ds, 10)

	visible := vl.RenderVisible()
	if len(visible) != 5 {
		t.Errorf("expected 5 visible items for 5 total rows, got %d", len(visible))
	}
}
