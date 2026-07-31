package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kenth/halptask/config"
	"github.com/kenth/halptask/model"
)

func TestNormalModeTaskShortcuts(t *testing.T) {
	cfg := config.DefaultConfig()
	tree := model.NewTree()
	item := tree.InsertBelow("", "Test Bullet")

	app := AppModel{
		Config:    cfg,
		Tree:      tree,
		Mode:      ModeNormal,
		WhichKey:  NewWhichKeyModel(),
		QuickHelp: NewQuickHelp(),
		TreeView:  NewTreeView(),
		StatusBar: NewStatusBar(),
		HelpModal: NewHelpModal(),
	}
	app.ensureValidCursor()
	app.SelectedID = item.ID

	// Test single key 't': Instant toggle bullet -> task (Todo)
	m, _ := app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	app = m.(AppModel)
	updatedItem := app.Tree.FindItem(item.ID)
	if !updatedItem.IsTask || updatedItem.Status != model.StatusTodo {
		t.Fatalf("expected item to be Todo task after single 't', got isTask=%v status=%s", updatedItem.IsTask, updatedItem.Status)
	}

	// Test single key 't' again: Cycle to InProgress
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	app = m.(AppModel)
	updatedItem = app.Tree.FindItem(item.ID)
	if updatedItem.Status != model.StatusInProgress {
		t.Fatalf("expected item status InProgress after second 't', got %s", updatedItem.Status)
	}

	// Test single key 't' third time: Cycle to Done
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	app = m.(AppModel)
	updatedItem = app.Tree.FindItem(item.ID)
	if updatedItem.Status != model.StatusDone {
		t.Fatalf("expected item status Done after third 't', got %s", updatedItem.Status)
	}

	// Test Leader key '<space> t s': Mark Todo
	app.tryExecuteKeyBinding([]string{" ", "t", "s"})
	updatedItem = app.Tree.FindItem(item.ID)
	if updatedItem.Status != model.StatusTodo {
		t.Fatalf("expected item status Todo after Leader '<space> t s', got %s", updatedItem.Status)
	}

	// Test Leader key '<space> t d': Mark Done
	app.tryExecuteKeyBinding([]string{" ", "t", "d"})
	updatedItem = app.Tree.FindItem(item.ID)
	if updatedItem.Status != model.StatusDone {
		t.Fatalf("expected item status Done after Leader '<space> t d', got %s", updatedItem.Status)
	}

	// Test Leader key '<space> t p': Mark In Progress
	app.tryExecuteKeyBinding([]string{" ", "t", "p"})
	updatedItem = app.Tree.FindItem(item.ID)
	if updatedItem.Status != model.StatusInProgress {
		t.Fatalf("expected item status InProgress after Leader '<space> t p', got %s", updatedItem.Status)
	}
}

func TestTier1AndTier2Shortcuts(t *testing.T) {
	cfg := config.DefaultConfig()
	tree := model.NewTree()
	parent := tree.InsertBelow("", "Parent Item")
	child := tree.AddChild(parent.ID, "Child Item")
	_ = child
	doneTask := tree.InsertBelow("", "Done Task")
	tree.SetStatus(doneTask.ID, model.StatusDone)

	app := AppModel{
		Config:    cfg,
		Tree:      tree,
		Mode:      ModeNormal,
		TextInput: textinput.New(),
		WhichKey:  NewWhichKeyModel(),
		QuickHelp: NewQuickHelp(),
		TreeView:  NewTreeView(),
		StatusBar: NewStatusBar(),
		HelpModal: NewHelpModal(),
	}
	app.ensureValidCursor()
	app.SelectedID = parent.ID

	// Test 'enter': Toggle fold
	m, _ := app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	_ = m
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyEnter})
	app = m.(AppModel)
	pItem := app.Tree.FindItem(parent.ID)
	if !pItem.Folded {
		t.Fatalf("expected parent item to be folded after pressing enter")
	}

	// Test 'c': Clear text and enter insert mode
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	app = m.(AppModel)
	if app.Mode != ModeInsert || app.TextInput.Value() != "" {
		t.Fatalf("expected insert mode with empty text after 'c', got mode=%v val=%q", app.Mode, app.TextInput.Value())
	}
	app.Mode = ModeNormal // reset back to normal

	// Test 'fc': Toggle hide completed
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	app = m.(AppModel)
	if !app.HideCompleted {
		t.Fatalf("expected HideCompleted to be true after 'fc'")
	}

	// Test 'ff': Zoom into focused subtree
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	app = m.(AppModel)
	if app.ZoomedID != parent.ID {
		t.Fatalf("expected ZoomedID to be parent.ID after 'ff', got %q", app.ZoomedID)
	}

	// Test 'da': Delete all done tasks
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = m.(AppModel)
	if app.Tree.FindItem(doneTask.ID) != nil {
		t.Fatalf("expected done task to be deleted after 'da'")
	}
}

func TestQuickHelpInView(t *testing.T) {
	cfg := config.DefaultConfig()
	tree := model.NewTree()
	tree.InsertBelow("", "Test Bullet")

	app := AppModel{
		Config:    cfg,
		Storage:   model.NewStorage("test.md", false),
		Tree:      tree,
		Mode:      ModeNormal,
		WhichKey:  NewWhichKeyModel(),
		QuickHelp: NewQuickHelp(),
		TreeView:  NewTreeView(),
		StatusBar: NewStatusBar(),
		HelpModal: NewHelpModal(),
		Width:     80,
		Height:    24,
	}

	viewStr := app.View()
	if !strings.Contains(viewStr, "<space>") || !strings.Contains(viewStr, "leader") {
		t.Fatalf("expected app View output to render quick help bar containing leader key shortcuts")
	}
}

func TestOpenAndChildShortcuts(t *testing.T) {
	cfg := config.DefaultConfig()
	tree := model.NewTree()
	item := tree.InsertBelow("", "Root Bullet 1")

	app := AppModel{
		Config:    cfg,
		Tree:      tree,
		Mode:      ModeNormal,
		TextInput: textinput.New(),
		WhichKey:  NewWhichKeyModel(),
		QuickHelp: NewQuickHelp(),
		TreeView:  NewTreeView(),
		StatusBar: NewStatusBar(),
		HelpModal: NewHelpModal(),
	}
	app.ensureValidCursor()
	app.SelectedID = item.ID

	// Test 'oc': Add child bullet
	m, _ := app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	app = m.(AppModel)

	if app.Mode != ModeInsert {
		t.Fatalf("expected insert mode after 'oc', got mode=%v", app.Mode)
	}
	rootItem := app.Tree.FindItem(item.ID)
	if len(rootItem.Children) != 1 {
		t.Fatalf("expected rootItem to have 1 child after 'oc', got %d", len(rootItem.Children))
	}

	// Test 'oo': New bullet below (sibling) on a bullet without children
	item2 := app.Tree.InsertBelow("", "Root Bullet 2")
	app.Mode = ModeNormal
	app.SelectedID = item2.ID

	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)

	if app.Mode != ModeInsert {
		t.Fatalf("expected insert mode after 'oo', got mode=%v", app.Mode)
	}
	if len(app.Tree.Roots) != 3 {
		t.Fatalf("expected 3 root items after 'oo', got %d", len(app.Tree.Roots))
	}
}

