# HalpTask Documentation Hub 📚

Welcome to the HalpTask documentation!

## 📖 Documentation Index

1. [Keybindings Cheatsheet](cheatsheet.md): Complete reference for leader key commands, direct Vim navigation, task state management, and folding.
2. [Agent Developer Guide](../AGENT.md): Complete guide for AI agents onboarding and working on HalpTask codebase.
3. **Architecture & Design**:
   - `config/`: Configuration loader and YAML parser (`~/.config/halptask/config.yaml`).
   - `model/`: Hierarchical tree model, Markdown parser/serializer, and AES-256-GCM encryption engine.
   - `ui/`: Lip Gloss styled components (WhichKey popup, TreeView, StatusBar, HelpModal) and Bubble Tea application lifecycle engine.

---

## 🛠️ Configuration & Storage

### Default Configuration Location
`~/.config/halptask/config.yaml`

```yaml
auto_save: true
data_file: ~/.config/halptask/data.txt
encrypted: false
indent_spaces: 2
leader_key: " "
show_which_key: true
theme: default
```

**Auto-Save & Encrypted Files**: 
When `auto_save: true`, HalpTask will automatically save all tree state mutations in the background. If you open or create an encrypted file but haven't provided a passphrase yet, auto-save will gracefully pause until you enter your passphrase.

### Data Storage Format
HalpTask saves data in standard human-readable Markdown format:

```markdown
- Welcome to HalpTask!
  - [ ] Todo item
  - [~] In progress item
  - [x] Completed item <!-- fold -->
    - Nested detail line
```
