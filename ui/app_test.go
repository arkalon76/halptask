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

func TestTagPickerAutoSave(t *testing.T) {
	tmpDir := t.TempDir()
	dataFile := tmpDir + "/tasks.txt"

	cfg := config.DefaultConfig()
	cfg.AutoSave = true
	cfg.DataFile = dataFile

	storage := model.NewStorage(dataFile, false)
	tree := model.NewTree()
	item := tree.InsertBelow("", "Build feature")

	app := AppModel{
		Config:    cfg,
		Storage:   storage,
		Tree:      tree,
		Mode:      ModeNormal,
		TagModal:  NewTagModal(cfg.Tags),
		WhichKey:  NewWhichKeyModel(),
		QuickHelp: NewQuickHelp(),
		TreeView:  NewTreeView(),
		StatusBar: NewStatusBar(),
		HelpModal: NewHelpModal(),
	}
	app.ensureValidCursor()
	app.SelectedID = item.ID

	// 1. Open TagPicker via 'T'
	m, _ := app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})
	app = m.(AppModel)
	if app.Mode != ModeTagPicker {
		t.Fatalf("Expected AppMode to be ModeTagPicker after pressing 'T', got %v", app.Mode)
	}

	// 2. Toggle tag 1 ("bug") in TagModal
	m, _ = app.Update(tea.KeyMsg{Type: tea.KeyEnter}) // toggle first tag
	app = m.(AppModel)

	targetItem := app.Tree.FindItem(item.ID)
	if !targetItem.HasDirectTag("bug") {
		t.Fatalf("Expected target item to have direct tag 'bug'")
	}

	// 3. Close TagPicker modal via 'esc'
	m, _ = app.Update(tea.KeyMsg{Type: tea.KeyEsc})
	app = m.(AppModel)

	if app.Mode != ModeNormal {
		t.Fatalf("Expected mode to reset to ModeNormal after esc, got %v", app.Mode)
	}

	// 4. Verify data file on disk contains the tag #bug!
	loadedTree, err := storage.Load("")
	if err != nil {
		t.Fatalf("Failed to load saved data file: %v", err)
	}
	if len(loadedTree.Roots) == 0 || !loadedTree.Roots[0].HasDirectTag("bug") {
		t.Fatalf("Saved file on disk missing tag 'bug'")
	}
}

func TestDefaultItemTypeCreation(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DefaultItemType = "task" // Set default creation to task

	tree := model.NewTree()
	item := tree.InsertBelow("", "Root Bullet")

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

	// Create new item below using 'oo'
	m, _ := app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)

	visible := app.getVisibleItems()
	if len(visible) != 2 {
		t.Fatalf("expected 2 visible items, got %d", len(visible))
	}

	createdItem := visible[1].Item
	if !createdItem.IsTask || createdItem.Status != model.StatusTodo {
		t.Fatalf("expected newly created item to be a Todo task, got isTask=%v status=%s", createdItem.IsTask, createdItem.Status)
	}

	// Switch default item type back to "bullet"
	app.Config.DefaultItemType = "bullet"
	app.Mode = ModeNormal

	// Create new item below using 'oo'
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)
	m, _ = app.updateNormal(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	app = m.(AppModel)

	visible = app.getVisibleItems()
	if len(visible) != 3 {
		t.Fatalf("expected 3 visible items, got %d", len(visible))
	}

	createdBullet := visible[2].Item
	if createdBullet.IsTask {
		t.Fatalf("expected newly created item to be a bullet point (isTask=false), got isTask=true")
	}
}

func TestToggleDefaultItemType(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DefaultItemType = "bullet"

	app := AppModel{
		Config:    cfg,
		Tree:      model.NewTree(),
		Mode:      ModeNormal,
		WhichKey:  NewWhichKeyModel(),
		QuickHelp: NewQuickHelp(),
		TreeView:  NewTreeView(),
		StatusBar: NewStatusBar(),
		HelpModal: NewHelpModal(),
	}

	// Toggle via Leader key '<space> t D'
	app.tryExecuteKeyBinding([]string{" ", "t", "D"})
	if app.Config.DefaultItemType != "task" {
		t.Fatalf("expected DefaultItemType to be 'task' after toggle, got %s", app.Config.DefaultItemType)
	}
	if !strings.Contains(app.StatusMsg, "Task") {
		t.Fatalf("expected status message to mention Task, got %q", app.StatusMsg)
	}

	// Toggle again
	app.tryExecuteKeyBinding([]string{" ", "t", "D"})
	if app.Config.DefaultItemType != "bullet" {
		t.Fatalf("expected DefaultItemType to be 'bullet' after second toggle, got %s", app.Config.DefaultItemType)
	}
	if !strings.Contains(app.StatusMsg, "Bullet") {
		t.Fatalf("expected status message to mention Bullet, got %q", app.StatusMsg)
	}
}

