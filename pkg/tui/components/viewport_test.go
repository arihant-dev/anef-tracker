package components_test

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/tui/components"
	"strings"
	"testing"
)

func TestViewportScrolling(t *testing.T) {
	vp := components.NewViewport(80, 5)

	var lines []string
	for i := 1; i <= 20; i++ {
		lines = append(lines, fmt.Sprintf("Line %d", i))
	}
	vp.SetContent(lines)

	if vp.Offset != 0 {
		t.Errorf("expected initial offset 0, got %d", vp.Offset)
	}

	vp.ScrollDown(3)
	if vp.Offset != 3 {
		t.Errorf("expected offset 3 after ScrollDown(3), got %d", vp.Offset)
	}

	vp.ScrollUp(2)
	if vp.Offset != 1 {
		t.Errorf("expected offset 1 after ScrollUp(2), got %d", vp.Offset)
	}
}

func TestViewportPageDown(t *testing.T) {
	vp := components.NewViewport(80, 5)

	var lines []string
	for i := 1; i <= 20; i++ {
		lines = append(lines, fmt.Sprintf("Line %d", i))
	}
	vp.SetContent(lines)

	vp.PageDown()
	if vp.Offset != 5 {
		t.Errorf("expected offset 5 after PageDown(), got %d", vp.Offset)
	}

	vp.PageUp()
	if vp.Offset != 0 {
		t.Errorf("expected offset 0 after PageUp(), got %d", vp.Offset)
	}
}

func TestViewportDoesNotOverflow(t *testing.T) {
	vp := components.NewViewport(80, 5)

	var lines []string
	for i := 1; i <= 10; i++ {
		lines = append(lines, fmt.Sprintf("Line %d", i))
	}
	vp.SetContent(lines)

	vp.ScrollDown(100)
	if vp.Offset != 5 { // maxOffset = 10 - 5 = 5
		t.Errorf("expected max offset 5, got %d", vp.Offset)
	}

	rendered := vp.Render()
	renderedLines := strings.Split(rendered, "\n")
	if len(renderedLines) != 5 {
		t.Errorf("expected exactly 5 rendered lines, got %d", len(renderedLines))
	}
}

func TestHeaderAlwaysVisible(t *testing.T) {
	vp := components.NewViewport(80, 5)
	vp.SetContent([]string{"Line 1", "Line 2", "Line 3"})

	rendered := vp.Render()
	lines := strings.Split(rendered, "\n")
	if len(lines) != 5 {
		t.Errorf("expected fixed height of 5 lines, got %d", len(lines))
	}
}
