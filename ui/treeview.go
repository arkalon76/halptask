package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kenth/halptask/model"
)

type TreeView struct {
	Width        int
	Height       int
	SelectedID   string
	SearchQuery  string
	MatchedIDs   map[string]bool
	IndentSpaces int
}

func NewTreeView() TreeView {
	return TreeView{
		Width:        80,
		Height:       20,
		IndentSpaces: 2,
		MatchedIDs:   make(map[string]bool),
	}
}

func (tv *TreeView) Render(visible []model.VisibleItem, cursorIndex int, scrollOffset int) string {
	if len(visible) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Italic(true).
			Padding(2, 4)
		return emptyStyle.Render("No bullets yet. Press 'o' or '<space> b n' to create a bullet point!")
	}

	// Styles
	cursorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7aa2f7"))

	indentGuideStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3b4261"))

	foldIconStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e0af68"))

	// Status styles
	todoBoxStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#787c99")) // Gray empty [ ]

	inProgressBoxStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff9e64")) // Orange [~]

	doneBoxStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9ece6a")) // Green [x]

	// Text styles
	normalTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c0caf5"))

	doneTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89")).
		Strikethrough(true)

	selectedRowStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#2e3c64")).
		Bold(true)

	searchMatchStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#e0af68")).
		Foreground(lipgloss.Color("#1a1b26")).
		Bold(true)

	bulletStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7dcfff"))

	var lines []string

	// Determine render window bounds based on Height
	maxLines := tv.Height
	if maxLines <= 0 {
		maxLines = 20
	}

	start := scrollOffset
	end := scrollOffset + maxLines
	if end > len(visible) {
		end = len(visible)
	}

	for i := start; i < end; i++ {
		v := visible[i]
		item := v.Item
		isSelected := (i == cursorIndex)

		// Indentation prefix
		indentStr := ""
		if v.Depth > 0 {
			parts := make([]string, v.Depth)
			for d := 0; d < v.Depth; d++ {
				parts[d] = indentGuideStyle.Render("│ ")
			}
			indentStr = strings.Join(parts, "")
		}

		// Cursor prefix
		cursorStr := "  "
		if isSelected {
			cursorStr = cursorStyle.Render("❯ ")
		}

		// Fold icon / Bullet prefix
		var prefix string
		if v.HasChildren {
			if item.Folded {
				childCount := len(item.Children)
				prefix = foldIconStyle.Render(fmt.Sprintf("▶ [%d] ", childCount))
			} else {
				prefix = foldIconStyle.Render("▼ ")
			}
		} else {
			prefix = bulletStyle.Render("• ")
		}

		// Status / Checkbox prefix
		var statusBox string
		if item.IsTask {
			switch item.Status {
			case model.StatusDone:
				statusBox = doneBoxStyle.Render("[x] ")
			case model.StatusInProgress:
				statusBox = inProgressBoxStyle.Render("[~] ")
			case model.StatusTodo:
				statusBox = todoBoxStyle.Render("[ ] ")
			default:
				statusBox = todoBoxStyle.Render("[ ] ")
			}
		}

		// Text rendering
		var formattedText string
		if item.IsTask && item.Status == model.StatusDone {
			formattedText = doneTextStyle.Render(item.Text)
		} else {
			formattedText = normalTextStyle.Render(item.Text)
		}

		// Search highlight if searching
		if tv.SearchQuery != "" && tv.MatchedIDs[item.ID] {
			formattedText = searchMatchStyle.Render(item.Text)
		}

		lineContent := fmt.Sprintf("%s%s%s%s%s", cursorStr, indentStr, prefix, statusBox, formattedText)

		if isSelected {
			lineContent = selectedRowStyle.Render(lineContent)
		}

		lines = append(lines, lineContent)
	}

	return strings.Join(lines, "\n")
}
