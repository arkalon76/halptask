package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type WhichKeyModel struct {
	Active     bool
	PrefixKeys []string
	Width      int
}

func NewWhichKeyModel() WhichKeyModel {
	return WhichKeyModel{
		Active:     false,
		PrefixKeys: []string{},
		Width:      80,
	}
}

type WhichKeyOption struct {
	Key     string
	Label   string
	IsGroup bool
}

func (wk *WhichKeyModel) GetOptions(allBindings []KeyBinding) (string, []WhichKeyOption) {
	prefix := wk.PrefixKeys
	title := "WhichKey"

	if len(prefix) > 0 {
		displayPrefix := make([]string, len(prefix))
		for i, k := range prefix {
			if k == " " {
				displayPrefix[i] = "Space"
			} else {
				displayPrefix[i] = k
			}
		}
		title = fmt.Sprintf(" %s ", strings.Join(displayPrefix, " > "))
	}

	optionsMap := make(map[string]WhichKeyOption)

	for _, b := range allBindings {
		if len(b.Keys) <= len(prefix) {
			continue
		}

		// Check if b.Keys matches prefix so far
		match := true
		for i := 0; i < len(prefix); i++ {
			if b.Keys[i] != prefix[i] {
				match = false
				break
			}
		}
		if !match {
			continue
		}

		nextKey := b.Keys[len(prefix)]
		nextKeyDisplay := nextKey
		if nextKey == " " {
			nextKeyDisplay = "<space>"
		}

		if len(b.Keys) > len(prefix)+1 {
			// It's a prefix group (e.g. 'b' -> bullets)
			groupName := b.Category
			if nextKey == "b" {
				groupName = "+bullets"
			} else if nextKey == "t" {
				groupName = "+tasks"
			} else if nextKey == "z" {
				groupName = "+folds"
			} else if nextKey == "e" {
				groupName = "+encrypt"
			} else if nextKey == "f" {
				groupName = "+file"
			} else {
				groupName = "+" + strings.ToLower(b.Category)
			}

			optionsMap[nextKey] = WhichKeyOption{
				Key:     nextKeyDisplay,
				Label:   groupName,
				IsGroup: true,
			}
		} else {
			// Direct command
			optionsMap[nextKey] = WhichKeyOption{
				Key:     nextKeyDisplay,
				Label:   b.Label,
				IsGroup: false,
			}
		}
	}

	var keys []string
	for k := range optionsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var options []WhichKeyOption
	for _, k := range keys {
		options = append(options, optionsMap[k])
	}

	return title, options
}

func (wk *WhichKeyModel) Render(allBindings []KeyBinding, width int) string {
	if !wk.Active {
		return ""
	}
	if width <= 0 {
		width = 80
	}

	titleText, options := wk.GetOptions(allBindings)
	if len(options) == 0 {
		return ""
	}

	// Lipgloss styles for TokioNight theme
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7aa2f7")).
		Padding(0, 1).
		Width(width - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1b26")).
		Background(lipgloss.Color("#7aa2f7")).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#e0af68")).
		Width(6)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a9b1d6"))

	groupStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#bb9af7"))

	var items []string
	for _, opt := range options {
		kStr := keyStyle.Render(opt.Key)
		var lStr string
		if opt.IsGroup {
			lStr = groupStyle.Render(opt.Label)
		} else {
			lStr = labelStyle.Render(opt.Label)
		}
		items = append(items, fmt.Sprintf("%s %s", kStr, lStr))
	}

	// Format into columns (3 or 4 columns)
	cols := 3
	if width > 100 {
		cols = 4
	}

	var rows []string
	for i := 0; i < len(items); i += cols {
		end := i + cols
		if end > len(items) {
			end = len(items)
		}

		var rowParts []string
		for j := i; j < end; j++ {
			// pad each column to ~ 24 chars
			part := lipgloss.NewStyle().Width((width - 8) / cols).Render(items[j])
			rowParts = append(rowParts, part)
		}
		rows = append(rows, strings.Join(rowParts, ""))
	}

	content := strings.Join(rows, "\n")
	header := titleStyle.Render(titleText)

	fullPopup := fmt.Sprintf("%s\n%s", header, content)
	return borderStyle.Render(fullPopup)
}
