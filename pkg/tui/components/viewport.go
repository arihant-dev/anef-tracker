package components

import (
	"fmt"
	"strings"
)

type Viewport struct {
	Width       int
	Height      int
	Offset      int
	Content     []string
	TotalItems  int
	CurrentPage int
	TotalPages  int
}

func NewViewport(width, height int) *Viewport {
	if height < 1 {
		height = 10
	}
	if width < 1 {
		width = 80
	}
	return &Viewport{
		Width:  width,
		Height: height,
	}
}

func (v *Viewport) SetContent(lines []string) {
	v.Content = lines
	v.UpdatePagination()
	if v.Offset > len(v.Content)-v.Height {
		v.Offset = len(v.Content) - v.Height
	}
	if v.Offset < 0 {
		v.Offset = 0
	}
}

func (v *Viewport) UpdatePagination() {
	if v.Height <= 0 {
		v.TotalPages = 1
		v.CurrentPage = 1
		return
	}

	totalLines := len(v.Content)
	if totalLines == 0 {
		v.TotalPages = 1
		v.CurrentPage = 1
		return
	}

	v.TotalPages = (totalLines + v.Height - 1) / v.Height
	v.CurrentPage = (v.Offset / v.Height) + 1
	if v.CurrentPage > v.TotalPages {
		v.CurrentPage = v.TotalPages
	}
}

func (v *Viewport) ScrollDown(lines int) {
	v.Offset += lines
	maxOffset := len(v.Content) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if v.Offset > maxOffset {
		v.Offset = maxOffset
	}
	v.UpdatePagination()
}

func (v *Viewport) ScrollUp(lines int) {
	v.Offset -= lines
	if v.Offset < 0 {
		v.Offset = 0
	}
	v.UpdatePagination()
}

func (v *Viewport) PageDown() {
	v.ScrollDown(v.Height)
}

func (v *Viewport) PageUp() {
	v.ScrollUp(v.Height)
}

func (v *Viewport) ScrollToTop() {
	v.Offset = 0
	v.UpdatePagination()
}

func (v *Viewport) ScrollToBottom() {
	maxOffset := len(v.Content) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	v.Offset = maxOffset
	v.UpdatePagination()
}

func (v Viewport) Render() string {
	if len(v.Content) == 0 {
		return ""
	}

	start := v.Offset
	if start >= len(v.Content) {
		start = len(v.Content) - 1
	}
	if start < 0 {
		start = 0
	}

	end := start + v.Height
	if end > len(v.Content) {
		end = len(v.Content)
	}

	visibleLines := v.Content[start:end]

	// Pad empty lines if content is shorter than height
	for len(visibleLines) < v.Height {
		visibleLines = append(visibleLines, "")
	}

	return strings.Join(visibleLines, "\n")
}

func (v Viewport) FormatPageIndicator(itemLabel string) string {
	totalLines := len(v.Content)
	if totalLines == 0 {
		return fmt.Sprintf("0 %s", itemLabel)
	}

	startItem := v.Offset + 1
	endItem := v.Offset + v.Height
	if endItem > totalLines {
		endItem = totalLines
	}

	return fmt.Sprintf("Showing %d-%d of %d %s (Page %d/%d)", startItem, endItem, totalLines, itemLabel, v.CurrentPage, v.TotalPages)
}
