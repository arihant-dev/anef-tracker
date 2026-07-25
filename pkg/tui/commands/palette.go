package commands

import (
	"fmt"
	"strings"
)

type CommandPalette struct {
	Active        bool
	Query         string
	SelectedIndex int
	Filtered      []CommandItem
}

func NewCommandPalette() *CommandPalette {
	cp := &CommandPalette{
		Filtered: Catalog,
	}
	return cp
}

func (cp *CommandPalette) UpdateQuery(q string) {
	cp.Query = q
	if q == "" {
		cp.Filtered = Catalog
		cp.SelectedIndex = 0
		return
	}

	var match []CommandItem
	queryLower := strings.ToLower(q)
	for _, item := range Catalog {
		if strings.Contains(strings.ToLower(item.Key), queryLower) ||
			strings.Contains(strings.ToLower(item.Title), queryLower) ||
			strings.Contains(strings.ToLower(item.Description), queryLower) {
			match = append(match, item)
		}
	}
	cp.Filtered = match
	cp.SelectedIndex = 0
}

func (cp *CommandPalette) SelectNext() {
	if len(cp.Filtered) == 0 {
		return
	}
	cp.SelectedIndex = (cp.SelectedIndex + 1) % len(cp.Filtered)
}

func (cp *CommandPalette) SelectPrev() {
	if len(cp.Filtered) == 0 {
		return
	}
	cp.SelectedIndex = (cp.SelectedIndex + len(cp.Filtered) - 1) % len(cp.Filtered)
}

func (cp CommandPalette) Render(width int) string {
	var sb strings.Builder
	sb.WriteString("=== GLOBAL COMMAND PALETTE (Ctrl+P) ===\n")
	sb.WriteString(fmt.Sprintf("> %s_\n", cp.Query))
	sb.WriteString(strings.Repeat("-", width-6))
	sb.WriteString("\n")

	if len(cp.Filtered) == 0 {
		sb.WriteString("No matching commands found.\n")
		return sb.String()
	}

	for i, item := range cp.Filtered {
		prefix := "  "
		if i == cp.SelectedIndex {
			prefix = " >"
		}
		sb.WriteString(fmt.Sprintf("%s [%s] %s — %s\n", prefix, item.Key, item.Title, item.Description))
	}

	return sb.String()
}
