package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kenth/halptask/config"
	"github.com/kenth/halptask/model"
)

type PromptType int

const (
	PromptPassphraseLoad PromptType = iota
	PromptPassphraseSet
)

type AppModel struct {
	Config       *config.Config
	Storage      *model.Storage
	Tree         *model.Tree
	UndoStack    []*model.Tree
	RedoStack    []*model.Tree

	Mode         AppMode
	CursorIndex  int
	ScrollOffset int
	SelectedID   string

	TextInput    textinput.Model
	SearchInput  textinput.Model
	PromptInput  textinput.Model
	PromptType   PromptType

	WhichKey     WhichKeyModel
	QuickHelp    QuickHelp
	TreeView     TreeView
	StatusBar    StatusBar
	HelpModal    HelpModal

	Passphrase    string
	StatusMsg     string
	Width         int
	Height        int

	ZoomedID      string
	HideCompleted bool

	// Key sequence buffer for multi-keystroke commands (e.g. "gg", "dd", "zc")
	KeyBuffer string
}

func (m *AppModel) getVisibleItems() []model.VisibleItem {
	return m.Tree.FlattenVisibleFiltered(m.ZoomedID, m.HideCompleted)
}

func InitialModel(cfg *config.Config, storage *model.Storage) (AppModel, tea.Cmd) {
	ti := textinput.New()
	ti.Prompt = "✏️  "
	ti.Placeholder = "Enter bullet text..."
	ti.CharLimit = 256

	si := textinput.New()
	si.Prompt = "🔍 "
	si.Placeholder = "Search bullets..."

	pi := textinput.New()
	pi.Prompt = "🔑 Passphrase: "
	pi.EchoMode = textinput.EchoPassword
	pi.EchoCharacter = '•'

	m := AppModel{
		Config:       cfg,
		Storage:      storage,
		Tree:         model.NewTree(),
		UndoStack:    []*model.Tree{},
		RedoStack:    []*model.Tree{},
		Mode:         ModeNormal,
		CursorIndex:  0,
		ScrollOffset: 0,
		TextInput:    ti,
		SearchInput:  si,
		PromptInput:  pi,
		WhichKey:     NewWhichKeyModel(),
		QuickHelp:    NewQuickHelp(),
		TreeView:     NewTreeView(),
		StatusBar:    NewStatusBar(),
		HelpModal:    NewHelpModal(),
		Width:        80,
		Height:       24,
	}

	// Check if target file is encrypted
	isEncrypted, err := model.IsEncryptedFile(storage.FilePath)
	if err == nil && isEncrypted {
		m.Mode = ModePrompt
		m.PromptType = PromptPassphraseLoad
		m.PromptInput.Focus()
		return m, textinput.Blink
	}

	// Otherwise load plain text file
	tree, err := storage.Load("")
	if err == nil {
		m.Tree = tree
		m.ensureValidCursor()
	}

	return m, nil
}

func (m *AppModel) pushUndo() {
	m.UndoStack = append(m.UndoStack, m.Tree.Clone())
	m.RedoStack = nil // Clear redo stack on new change
}

func (m *AppModel) undo() {
	if len(m.UndoStack) == 0 {
		m.StatusMsg = "Already at oldest change"
		return
	}
	m.RedoStack = append(m.RedoStack, m.Tree.Clone())
	m.Tree = m.UndoStack[len(m.UndoStack)-1]
	m.UndoStack = m.UndoStack[:len(m.UndoStack)-1]
	m.ensureValidCursor()
	m.StatusMsg = "Undo"
}

func (m *AppModel) redo() {
	if len(m.RedoStack) == 0 {
		m.StatusMsg = "Already at newest change"
		return
	}
	m.UndoStack = append(m.UndoStack, m.Tree.Clone())
	m.Tree = m.RedoStack[len(m.RedoStack)-1]
	m.RedoStack = m.RedoStack[:len(m.RedoStack)-1]
	m.ensureValidCursor()
	m.StatusMsg = "Redo"
}

func (m *AppModel) ensureValidCursor() {
	visible := m.getVisibleItems()
	if len(visible) == 0 {
		m.CursorIndex = 0
		m.SelectedID = ""
		return
	}
	if m.CursorIndex < 0 {
		m.CursorIndex = 0
	}
	if m.CursorIndex >= len(visible) {
		m.CursorIndex = len(visible) - 1
	}
	m.SelectedID = visible[m.CursorIndex].Item.ID

	// Adjust scroll window
	maxVisibleLines := m.Height - 5 // reserved for header/quickhelp/status/whichkey
	if maxVisibleLines < 5 {
		maxVisibleLines = 5
	}
	if m.CursorIndex < m.ScrollOffset {
		m.ScrollOffset = m.CursorIndex
	} else if m.CursorIndex >= m.ScrollOffset+maxVisibleLines {
		m.ScrollOffset = m.CursorIndex - maxVisibleLines + 1
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.TreeView.Width = msg.Width
		m.TreeView.Height = msg.Height - 6
		m.WhichKey.Width = msg.Width
		m.QuickHelp.Width = msg.Width
		m.StatusBar.Width = msg.Width
		m.HelpModal.Width = msg.Width
		m.HelpModal.Height = msg.Height

	case tea.KeyMsg:
		switch m.Mode {
		case ModeNormal:
			return m.updateNormal(msg)
		case ModeInsert:
			return m.updateInsert(msg)
		case ModeSearch:
			return m.updateSearch(msg)
		case ModePrompt:
			return m.updatePrompt(msg)
		case ModeHelp:
			if msg.String() == "esc" || msg.String() == "?" || msg.String() == "q" {
				m.Mode = ModeNormal
			}
			return m, nil
		}
	}

	return m, cmd
}

func (m AppModel) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	// Clear status message on key press
	m.StatusMsg = ""

	// Handle WhichKey popup navigation
	if m.WhichKey.Active || k == " " {
		if k == "esc" {
			m.WhichKey.Active = false
			m.WhichKey.PrefixKeys = nil
			m.KeyBuffer = ""
			return m, nil
		}

		if !m.WhichKey.Active && k == " " {
			m.WhichKey.Active = true
			m.WhichKey.PrefixKeys = []string{" "}
			return m, nil
		}

		// Keystroke while WhichKey active
		m.WhichKey.PrefixKeys = append(m.WhichKey.PrefixKeys, k)
		m.KeyBuffer = strings.Join(m.WhichKey.PrefixKeys, " ")

		actionExecuted := m.tryExecuteKeyBinding(m.WhichKey.PrefixKeys)
		if actionExecuted {
			m.WhichKey.Active = false
			m.WhichKey.PrefixKeys = nil
			m.KeyBuffer = ""
			return m, nil
		}

		// Check if prefix keys still match any known command
		title, options := m.WhichKey.GetOptions(GetAllKeyBindings())
		if len(options) == 0 && title == "WhichKey" {
			m.WhichKey.Active = false
			m.WhichKey.PrefixKeys = nil
			m.KeyBuffer = ""
		}
		return m, nil
	}

	// Handle single or double keystrokes in normal mode
	m.KeyBuffer += k

	// Check multi-key direct matches (e.g. "gg", "dd", "zc", "zo", "za", "zM", "zR", "ts")
	switch m.KeyBuffer {
	case "gg":
		m.CursorIndex = 0
		m.ensureValidCursor()
		m.KeyBuffer = ""
		return m, nil
	case "dd":
		if m.SelectedID != "" {
			m.pushUndo()
			nextID := m.Tree.Delete(m.SelectedID)
			m.ensureValidCursor()
			if nextID != "" {
				visible := m.getVisibleItems()
				for i, v := range visible {
					if v.Item.ID == nextID {
						m.CursorIndex = i
						break
					}
				}
			}
			m.ensureValidCursor()
			m.StatusMsg = "Deleted bullet"
		}
		m.KeyBuffer = ""
		return m, nil
	case "zc":
		if m.SelectedID != "" {
			m.Tree.Fold(m.SelectedID)
			m.ensureValidCursor()
		}
		m.KeyBuffer = ""
		return m, nil
	case "zo":
		if m.SelectedID != "" {
			m.Tree.Unfold(m.SelectedID)
			m.ensureValidCursor()
		}
		m.KeyBuffer = ""
		return m, nil
	case "za":
		if m.SelectedID != "" {
			m.Tree.ToggleFold(m.SelectedID)
			m.ensureValidCursor()
		}
		m.KeyBuffer = ""
		return m, nil
	case "zM":
		m.Tree.FoldAll()
		m.ensureValidCursor()
		m.KeyBuffer = ""
		return m, nil
	case "zR":
		m.Tree.UnfoldAll()
		m.ensureValidCursor()
		m.KeyBuffer = ""
		return m, nil
	case "ww":
		_ = m.saveFile()
		m.StatusMsg = "File saved"
		m.KeyBuffer = ""
		return m, nil
	case "fc":
		m.HideCompleted = !m.HideCompleted
		if m.HideCompleted {
			m.StatusMsg = "Hiding completed tasks [x]"
		} else {
			m.StatusMsg = "Showing all tasks"
		}
		m.ensureValidCursor()
		m.KeyBuffer = ""
		return m, nil
	case "da":
		m.pushUndo()
		count := m.Tree.DeleteCompleted()
		m.ensureValidCursor()
		m.StatusMsg = fmt.Sprintf("Cleared %d completed task(s)", count)
		m.KeyBuffer = ""
		return m, nil
	case "ff":
		if m.ZoomedID == "" && m.SelectedID != "" {
			m.ZoomedID = m.SelectedID
			item := m.Tree.FindItem(m.ZoomedID)
			if item != nil {
				m.StatusMsg = fmt.Sprintf("Zoomed in: %s", item.Text)
			}
		} else {
			m.ZoomedID = ""
			m.StatusMsg = "Unzoomed (full view)"
		}
		m.ensureValidCursor()
		m.KeyBuffer = ""
		return m, nil
	case "oo":
		m.pushUndo()
		newItem := m.Tree.InsertBelow(m.SelectedID, "")
		m.ensureValidCursor()
		visible := m.getVisibleItems()
		for i, v := range visible {
			if v.Item.ID == newItem.ID {
				m.CursorIndex = i
				break
			}
		}
		m.ensureValidCursor()
		m.Mode = ModeInsert
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.KeyBuffer = ""
		return m, textinput.Blink
	case "oc":
		m.pushUndo()
		newItem := m.Tree.AddChild(m.SelectedID, "")
		m.ensureValidCursor()
		visible := m.getVisibleItems()
		for i, v := range visible {
			if v.Item.ID == newItem.ID {
				m.CursorIndex = i
				break
			}
		}
		m.ensureValidCursor()
		m.Mode = ModeInsert
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.KeyBuffer = ""
		return m, textinput.Blink
	}

	// Single key handlers
	switch k {
	case "j", "down":
		m.CursorIndex++
		m.ensureValidCursor()
		m.KeyBuffer = ""
	case "k", "up":
		m.CursorIndex--
		m.ensureValidCursor()
		m.KeyBuffer = ""
	case "h", "left":
		if m.SelectedID != "" {
			item := m.Tree.FindItem(m.SelectedID)
			if item != nil && len(item.Children) > 0 && !item.Folded {
				m.Tree.Fold(m.SelectedID)
				m.ensureValidCursor()
			} else if item != nil && item.Parent != nil {
				// Jump to parent
				visible := m.getVisibleItems()
				for i, v := range visible {
					if v.Item.ID == item.Parent.ID {
						m.CursorIndex = i
						break
					}
				}
				m.ensureValidCursor()
			}
		}
		m.KeyBuffer = ""
	case "l", "right":
		if m.SelectedID != "" {
			item := m.Tree.FindItem(m.SelectedID)
			if item != nil && len(item.Children) > 0 && item.Folded {
				m.Tree.Unfold(m.SelectedID)
				m.ensureValidCursor()
			} else if item != nil && len(item.Children) > 0 {
				// Jump to first child
				visible := m.getVisibleItems()
				for i, v := range visible {
					if v.Item.ID == item.Children[0].ID {
						m.CursorIndex = i
						break
					}
				}
				m.ensureValidCursor()
			}
		}
		m.KeyBuffer = ""
	case "G":
		visible := m.getVisibleItems()
		if len(visible) > 0 {
			m.CursorIndex = len(visible) - 1
			m.ensureValidCursor()
		}
		m.KeyBuffer = ""
	case "O":
		m.pushUndo()
		newItem := m.Tree.InsertAbove(m.SelectedID, "")
		m.ensureValidCursor()
		visible := m.getVisibleItems()
		for i, v := range visible {
			if v.Item.ID == newItem.ID {
				m.CursorIndex = i
				break
			}
		}
		m.ensureValidCursor()
		m.Mode = ModeInsert
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.KeyBuffer = ""
		return m, textinput.Blink
	case "enter":
		if m.SelectedID != "" {
			m.Tree.ToggleFold(m.SelectedID)
			m.ensureValidCursor()
		}
		m.KeyBuffer = ""
	case "i", "a", "e":
		if m.SelectedID != "" {
			item := m.Tree.FindItem(m.SelectedID)
			if item != nil {
				m.pushUndo()
				m.Mode = ModeInsert
				m.TextInput.SetValue(item.Text)
				m.TextInput.Focus()
				m.KeyBuffer = ""
				return m, textinput.Blink
			}
		}
		m.KeyBuffer = ""
	case "c":
		if m.SelectedID != "" {
			item := m.Tree.FindItem(m.SelectedID)
			if item != nil {
				m.pushUndo()
				m.Mode = ModeInsert
				m.TextInput.SetValue("")
				m.TextInput.Focus()
				m.KeyBuffer = ""
				return m, textinput.Blink
			}
		}
		m.KeyBuffer = ""
	case "x":
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.Delete(m.SelectedID)
			m.ensureValidCursor()
			m.StatusMsg = "Deleted bullet"
		}
		m.KeyBuffer = ""
	case "t":
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.CycleStatus(m.SelectedID)
			m.ensureValidCursor()
		}
		m.KeyBuffer = ""
	case "tab":
		if m.SelectedID != "" {
			m.pushUndo()
			if m.Tree.Indent(m.SelectedID) {
				m.ensureValidCursor()
				m.StatusMsg = "Indented"
			}
		}
		m.KeyBuffer = ""
	case "shift+tab":
		if m.SelectedID != "" {
			m.pushUndo()
			if m.Tree.Unindent(m.SelectedID) {
				m.ensureValidCursor()
				m.StatusMsg = "Unindented"
			}
		}
		m.KeyBuffer = ""
	case "J":
		if m.SelectedID != "" {
			m.pushUndo()
			if m.Tree.MoveDown(m.SelectedID) {
				m.CursorIndex++
				m.ensureValidCursor()
			}
		}
		m.KeyBuffer = ""
	case "K":
		if m.SelectedID != "" {
			m.pushUndo()
			if m.Tree.MoveUp(m.SelectedID) {
				m.CursorIndex--
				m.ensureValidCursor()
			}
		}
		m.KeyBuffer = ""
	case "u":
		m.undo()
		m.KeyBuffer = ""
	case "ctrl+r":
		m.redo()
		m.KeyBuffer = ""
	case "/":
		m.Mode = ModeSearch
		m.SearchInput.Focus()
		m.KeyBuffer = ""
		return m, textinput.Blink
	case "?":
		m.Mode = ModeHelp
		m.KeyBuffer = ""
	case "q":
		_ = m.saveFile()
		return m, tea.Quit
	default:
		// If key buffer has accumulated keys but no match yet, check length
		if len(m.KeyBuffer) > 2 {
			m.KeyBuffer = ""
		}
	}

	return m, nil
}

func (m *AppModel) tryExecuteKeyBinding(keys []string) bool {
	keySeq := strings.Join(keys, " ")

	switch keySeq {
	case "  b n": // Space > Bullets > New below
		m.pushUndo()
		newItem := m.Tree.InsertBelow(m.SelectedID, "")
		m.ensureValidCursor()
		visible := m.getVisibleItems()
		for i, v := range visible {
			if v.Item.ID == newItem.ID {
				m.CursorIndex = i
				break
			}
		}
		m.ensureValidCursor()
		m.Mode = ModeInsert
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		return true
	case "  b N": // Space > Bullets > New above
		m.pushUndo()
		newItem := m.Tree.InsertAbove(m.SelectedID, "")
		m.ensureValidCursor()
		visible := m.getVisibleItems()
		for i, v := range visible {
			if v.Item.ID == newItem.ID {
				m.CursorIndex = i
				break
			}
		}
		m.ensureValidCursor()
		m.Mode = ModeInsert
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		return true
	case "  b c": // Space > Bullets > Add child
		if m.SelectedID != "" {
			m.pushUndo()
			newItem := m.Tree.AddChild(m.SelectedID, "")
			m.ensureValidCursor()
			visible := m.getVisibleItems()
			for i, v := range visible {
				if v.Item.ID == newItem.ID {
					m.CursorIndex = i
					break
				}
			}
			m.ensureValidCursor()
			m.Mode = ModeInsert
			m.TextInput.SetValue("")
			m.TextInput.Focus()
			return true
		}
	case "  b e": // Space > Bullets > Edit
		if m.SelectedID != "" {
			item := m.Tree.FindItem(m.SelectedID)
			if item != nil {
				m.pushUndo()
				m.Mode = ModeInsert
				m.TextInput.SetValue(item.Text)
				m.TextInput.Focus()
				return true
			}
		}
	case "  b d": // Space > Bullets > Delete
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.Delete(m.SelectedID)
			m.ensureValidCursor()
			m.StatusMsg = "Deleted bullet"
			return true
		}
	case "  b i": // Indent
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.Indent(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  b o": // Unindent
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.Unindent(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  b j": // Move down
		if m.SelectedID != "" {
			m.pushUndo()
			if m.Tree.MoveDown(m.SelectedID) {
				m.CursorIndex++
				m.ensureValidCursor()
			}
			return true
		}
	case "  b k": // Move up
		if m.SelectedID != "" {
			m.pushUndo()
			if m.Tree.MoveUp(m.SelectedID) {
				m.CursorIndex--
				m.ensureValidCursor()
			}
			return true
		}
	case "  b t": // Toggle task
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.ToggleTask(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  t t": // Toggle task
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.ToggleTask(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  t c": // Cycle task status
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.CycleStatus(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  t d": // Mark Done [x]
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.SetStatus(m.SelectedID, model.StatusDone)
			m.ensureValidCursor()
			m.StatusMsg = "Marked done [x]"
			return true
		}
	case "  t p": // Mark In Progress [~]
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.SetStatus(m.SelectedID, model.StatusInProgress)
			m.ensureValidCursor()
			m.StatusMsg = "Marked in-progress [~]"
			return true
		}
	case "  t s": // Mark Todo [ ]
		if m.SelectedID != "" {
			m.pushUndo()
			m.Tree.SetStatus(m.SelectedID, model.StatusTodo)
			m.ensureValidCursor()
			m.StatusMsg = "Marked todo [ ]"
			return true
		}
	case "  z c":
		if m.SelectedID != "" {
			m.Tree.Fold(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  z o":
		if m.SelectedID != "" {
			m.Tree.Unfold(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  z a":
		if m.SelectedID != "" {
			m.Tree.ToggleFold(m.SelectedID)
			m.ensureValidCursor()
			return true
		}
	case "  z M":
		m.Tree.FoldAll()
		m.ensureValidCursor()
		return true
	case "  z R":
		m.Tree.UnfoldAll()
		m.ensureValidCursor()
		return true
	case "  e e": // Toggle encryption
		m.Storage.Encrypted = !m.Storage.Encrypted
		if m.Storage.Encrypted && m.Passphrase == "" {
			m.Mode = ModePrompt
			m.PromptType = PromptPassphraseSet
			m.PromptInput.SetValue("")
			m.PromptInput.Focus()
		} else {
			_ = m.saveFile()
			if m.Storage.Encrypted {
				m.StatusMsg = "Encryption ENABLED 🔒"
			} else {
				m.StatusMsg = "Encryption DISABLED 🔓"
			}
		}
		return true
	case "  e p": // Change passphrase
		m.Mode = ModePrompt
		m.PromptType = PromptPassphraseSet
		m.PromptInput.SetValue("")
		m.PromptInput.Focus()
		return true
	case "  w", "  s":
		_ = m.saveFile()
		return true
	case "  /":
		m.Mode = ModeSearch
		m.SearchInput.Focus()
		return true
	case "  ?":
		m.Mode = ModeHelp
		return true
	case "  q":
		_ = m.saveFile()
		// Return true so updating stops, quit will be dispatched
		return true
	}

	return false
}

func (m AppModel) updateInsert(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "enter":
		text := strings.TrimSpace(m.TextInput.Value())
		if m.SelectedID != "" {
			item := m.Tree.FindItem(m.SelectedID)
			if item != nil {
				item.Text = text
			}
		}
		m.Mode = ModeNormal
		m.TextInput.Blur()
		_ = m.saveFile()
		return m, nil
	case "esc":
		m.Mode = ModeNormal
		m.TextInput.Blur()
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m AppModel) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc", "enter":
		m.Mode = ModeNormal
		m.SearchInput.Blur()
		return m, nil
	}

	m.SearchInput, cmd = m.SearchInput.Update(msg)
	query := m.SearchInput.Value()
	m.TreeView.SearchQuery = query

	if query != "" {
		matchedIDs := m.Tree.Search(query)
		m.TreeView.MatchedIDs = make(map[string]bool)
		for _, id := range matchedIDs {
			m.TreeView.MatchedIDs[id] = true
		}
		if len(matchedIDs) > 0 {
			visible := m.getVisibleItems()
			for i, v := range visible {
				if v.Item.ID == matchedIDs[0] {
					m.CursorIndex = i
					break
				}
			}
			m.ensureValidCursor()
		}
	} else {
		m.TreeView.MatchedIDs = make(map[string]bool)
	}

	return m, cmd
}

func (m AppModel) updatePrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "enter":
		val := m.PromptInput.Value()
		if m.PromptType == PromptPassphraseLoad {
			tree, err := m.Storage.Load(val)
			if err != nil {
				m.StatusMsg = "Incorrect passphrase! Try again."
				m.PromptInput.SetValue("")
				return m, textinput.Blink
			}
			m.Tree = tree
			m.Passphrase = val
			m.Storage.Encrypted = true
			m.Mode = ModeNormal
			m.ensureValidCursor()
			m.StatusMsg = "Decrypted & loaded successfully 🔓"
			return m, nil
		} else if m.PromptType == PromptPassphraseSet {
			m.Passphrase = val
			m.Storage.Encrypted = true
			m.Mode = ModeNormal
			_ = m.saveFile()
			m.StatusMsg = "Passphrase set & encrypted 🔒"
			return m, nil
		}
	case "esc":
		if m.PromptType == PromptPassphraseLoad {
			// Quit if user cancels load prompt
			return m, tea.Quit
		}
		m.Mode = ModeNormal
		return m, nil
	}

	m.PromptInput, cmd = m.PromptInput.Update(msg)
	return m, cmd
}

func (m *AppModel) saveFile() error {
	err := m.Storage.Save(m.Tree, m.Passphrase)
	if err != nil {
		m.StatusMsg = "Save error: " + err.Error()
		return err
	}
	m.StatusMsg = "Saved " + m.Storage.FilePath
	return nil
}

func (m AppModel) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading HalpTask..."
	}

	if m.Mode == ModePrompt {
		promptStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#f7768e")).
			Padding(1, 3).
			Width(60)

		title := "🔒 Encrypted File Passphrase"
		if m.PromptType == PromptPassphraseSet {
			title = "🔐 Set Encryption Passphrase"
		}

		content := fmt.Sprintf("%s\n\n%s", lipgloss.NewStyle().Bold(true).Render(title), m.PromptInput.View())
		if m.StatusMsg != "" {
			content += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e")).Render(m.StatusMsg)
		}
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, promptStyle.Render(content))
	}

	if m.Mode == ModeHelp {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, m.HelpModal.Render(m.Width, m.Height))
	}

	// Normal View Layout: Header / Tree View / Text Input / WhichKey / Status Bar
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1b26")).
		Background(lipgloss.Color("#7aa2f7")).
		Width(m.Width).
		Align(lipgloss.Center)

	header := headerStyle.Render(" HALPTASK  •  Bullet & Task Manager ")

	visible := m.getVisibleItems()
	treeContent := m.TreeView.Render(visible, m.CursorIndex, m.ScrollOffset)

	var midSection string
	if m.Mode == ModeInsert {
		insertBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9ece6a")).
			Padding(0, 1).
			Width(m.Width - 4).
			Render(m.TextInput.View())
		midSection = insertBox
	} else if m.Mode == ModeSearch {
		searchBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#e0af68")).
			Padding(0, 1).
			Width(m.Width - 4).
			Render(m.SearchInput.View())
		midSection = searchBox
	}

	var whichKeyStr string
	if m.WhichKey.Active {
		whichKeyStr = m.WhichKey.Render(GetAllKeyBindings(), m.Width)
	}

	stats := m.Tree.GetStats()
	filePath := ""
	isEncrypted := false
	if m.Storage != nil {
		filePath = m.Storage.FilePath
		isEncrypted = m.Storage.Encrypted
	}
	statusBarStr := m.StatusBar.Render(m.Mode, filePath, isEncrypted, stats, m.CursorIndex+1, len(visible), m.StatusMsg)

	var viewParts []string
	viewParts = append(viewParts, header)
	viewParts = append(viewParts, treeContent)

	if midSection != "" {
		viewParts = append(viewParts, midSection)
	}
	if whichKeyStr != "" {
		viewParts = append(viewParts, whichKeyStr)
	}
	viewParts = append(viewParts, m.QuickHelp.Render(m.Width))
	viewParts = append(viewParts, statusBarStr)

	return strings.Join(viewParts, "\n")
}
