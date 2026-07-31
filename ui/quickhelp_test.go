package ui

import (
	"strings"
	"testing"
)

func TestQuickHelpRender(t *testing.T) {
	qh := NewQuickHelp()
	rendered := qh.Render(80)

	if !strings.Contains(rendered, "<space>") {
		t.Errorf("expected rendered quick help to contain '<space>' leader shortcut")
	}
	if !strings.Contains(rendered, "j/k") {
		t.Errorf("expected rendered quick help to contain 'j/k' move shortcut")
	}
	if !strings.Contains(rendered, "o") {
		t.Errorf("expected rendered quick help to contain 'o' new shortcut")
	}
	if !strings.Contains(rendered, "tab") {
		t.Errorf("expected rendered quick help to contain 'tab' sub-bullet shortcut")
	}
	if !strings.Contains(rendered, "enter") {
		t.Errorf("expected rendered quick help to contain 'enter' fold shortcut")
	}
	if !strings.Contains(rendered, "i") {
		t.Errorf("expected rendered quick help to contain 'i' edit shortcut")
	}
	if !strings.Contains(rendered, "t") {
		t.Errorf("expected rendered quick help to contain 't' task shortcut")
	}
	if !strings.Contains(rendered, "?") {
		t.Errorf("expected rendered quick help to contain '?' help shortcut")
	}
}

func TestQuickHelpNarrowWidth(t *testing.T) {
	qh := NewQuickHelp()
	rendered := qh.Render(40)

	if !strings.Contains(rendered, "<space>") {
		t.Errorf("expected narrow rendered quick help to contain '<space>'")
	}
}
