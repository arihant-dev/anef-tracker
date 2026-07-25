package components

import (
	"fmt"
	"strings"
)

type DataSource interface {
	Count() int
	Row(index int) string
}

type SliceDataSource struct {
	Items []string
}

func (s *SliceDataSource) Count() int {
	return len(s.Items)
}

func (s *SliceDataSource) Row(index int) string {
	if index < 0 || index >= len(s.Items) {
		return ""
	}
	return s.Items[index]
}

type VirtualList struct {
	Offset     int
	Height     int
	Source     DataSource
	SelectedID int
}

func NewVirtualList(source DataSource, height int) *VirtualList {
	if height < 1 {
		height = 10
	}
	return &VirtualList{
		Source: source,
		Height: height,
	}
}

func (vl *VirtualList) SetSource(source DataSource) {
	vl.Source = source
	vl.ClampOffset()
}

func (vl *VirtualList) ClampOffset() {
	if vl.Source == nil {
		vl.Offset = 0
		return
	}
	total := vl.Source.Count()
	maxOffset := total - vl.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if vl.Offset > maxOffset {
		vl.Offset = maxOffset
	}
	if vl.Offset < 0 {
		vl.Offset = 0
	}
}

func (vl *VirtualList) ScrollDown(lines int) {
	vl.Offset += lines
	vl.ClampOffset()
}

func (vl *VirtualList) ScrollUp(lines int) {
	vl.Offset -= lines
	vl.ClampOffset()
}

func (vl *VirtualList) PageDown() {
	vl.ScrollDown(vl.Height)
}

func (vl *VirtualList) PageUp() {
	vl.ScrollUp(vl.Height)
}

func (vl VirtualList) RenderVisible() []string {
	if vl.Source == nil || vl.Source.Count() == 0 {
		return []string{}
	}

	total := vl.Source.Count()
	start := vl.Offset
	if start >= total {
		start = total - 1
	}
	if start < 0 {
		start = 0
	}

	end := start + vl.Height
	if end > total {
		end = total
	}

	var visible []string
	for i := start; i < end; i++ {
		visible = append(visible, vl.Source.Row(i))
	}

	return visible
}

func (vl VirtualList) RenderString() string {
	visible := vl.RenderVisible()
	return strings.Join(visible, "\n")
}

func (vl VirtualList) FormatPagination(itemLabel string) string {
	if vl.Source == nil || vl.Source.Count() == 0 {
		return fmt.Sprintf("0 %s", itemLabel)
	}

	total := vl.Source.Count()
	start := vl.Offset + 1
	end := vl.Offset + vl.Height
	if end > total {
		end = total
	}

	totalPages := (total + vl.Height - 1) / vl.Height
	currentPage := (vl.Offset / vl.Height) + 1
	if currentPage > totalPages {
		currentPage = totalPages
	}

	return fmt.Sprintf("Items %d-%d of %d %s (Page %d/%d)", start, end, total, itemLabel, currentPage, totalPages)
}
