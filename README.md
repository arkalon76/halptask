# HalpTask 🚀

A high-performance, keyboard-driven Terminal User Interface (TUI) bullet point and task manager written in **Go** and **Bubble Tea**, using well-known shortcuts.

![HalpTask TUI](https://img.shields.io/badge/TUI-Bubble%20Tea-purple)
![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 📚 Documentation & Cheatsheet

Detailed documentation and guides are available in the [`docs/`](docs/index.md) folder:

- **[Keybindings Cheatsheet](docs/cheatsheet.md)**: Full interactive Leader Key (`<space>`) map and Vim shortcuts.
- **[Documentation Index](docs/index.md)**: Architecture, data persistence, and configuration guide.

---

## ✨ Features

- **Vim Native Keybindings**: Native Vim navigation (`j`, `k`, `h`, `l`, `gg`, `G`, `oo`, `oc`, `O`, `dd`, `x`, `u`, `ctrl+r`, `tab`, `shift+tab`).
- **Leader Menu (`<space>`)**: Leader key popup window displaying all available shortcuts non-intrusively.
- **Dynamic WhichKey Popup**: Visual popup updates as you type key prefixes (e.g. `<space> b` for bullets, `<space> t` for tasks, `<space> z` for folds).
- **Bullet & Task Management**:
  - Convert any bullet point into a task with checkbox statuses.
  - **Todo**: `[ ]` (Gray empty box)
  - **In Progress**: `[~]` (Orange `~` indicator)
  - **Done**: `[x]` (Green `X` checkmark with faint strikethrough text styling)
- **Hierarchical Folding**:
  - Collapse (`zc`, `h`), Expand (`zo`, `l`), Toggle (`za`), Close All (`zM`), Open All (`zR`).
  - Child count badges for collapsed subtrees (`▶ [3]`).
- **Plain Text & Encryption**:
  - Data stored in clean human-readable Markdown format by default (`~/.config/halptask/data.txt`).
  - **AES-256-GCM + PBKDF2** encryption mode for secure task storage.
- **Cross Platform Support**:
  - Binaries built for **macOS** and **Linux** (`amd64` and `arm64`).
- **Customizable Configuration**:
  - Configured via `~/.config/halptask/config.yaml`.

---

## ⚡ Quick Start

### Build & Install

```bash
# Clone the repository
git clone https://github.com/kenth/halptask.git
cd halptask

# Build local binary
go build -o halptask .

# Run halptask
./halptask
```

### CLI Flags

```bash
./halptask -f ~/.config/halptask/my_tasks.txt    # Custom data file path
./halptask --encrypt                               # Force prompt for encryption
./halptask --version                               # Print version info
```

---

## ⌨️ Summary Keybindings Quick Reference

| Leader Key | Category | Command / Action |
|---|---|---|
| `<space> b n` | Bullets | New bullet below |
| `<space> b N` | Bullets | New bullet above |
| `<space> b c` | Bullets | Add child bullet |
| `<space> b e` | Bullets | Edit bullet text |
| `<space> b d` | Bullets | Delete bullet & subtree |
| `<space> b i` | Bullets | Indent bullet (demote) |
| `<space> b o` | Bullets | Unindent bullet (promote) |
| `<space> b j` | Bullets | Move bullet down |
| `<space> b k` | Bullets | Move bullet up |
| `<space> t t` | Tasks | Toggle bullet into task `[ ]` |
| `<space> t c` | Tasks | Cycle status (`Todo` ➜ `In Progress` ➜ `Done`) |
| `<space> t d` | Tasks | Mark Done `[x]` (Green X, strikethrough) |
| `<space> t p` | Tasks | Mark In Progress `[~]` (Orange ~) |
| `<space> t s` | Tasks | Mark Todo `[ ]` (Gray empty) |
| `<space> z c` | Folds | Close fold |
| `<space> z o` | Folds | Open fold |
| `<space> z a` | Folds | Toggle fold |
| `<space> z M` | Folds | Close all folds |
| `<space> z R` | Folds | Open all folds |
| `<space> e e` | Encrypt | Toggle encryption |
| `<space> e p` | Encrypt | Set / Change passphrase |
| `<space> w` | File | Save file |
| `<space> /` | Search | Search bullet points |
| `<space> ?` | Help | Show keymap cheat sheet modal |
| `<space> q` | Quit | Save and exit |

*For full details, view the [Keybindings Cheatsheet](docs/cheatsheet.md).*

---

## 🔒 Configuration File (`~/.config/halptask/config.yaml`)

```yaml
data_file: ~/.config/halptask/data.txt
encrypted: false
indent_spaces: 2
leader_key: " "
show_which_key: true
theme: default
```

---

## 🛠️ Cross-Compilation

To build binaries for all supported platforms:

```bash
# macOS (Apple Silicon & Intel)
GOOS=darwin GOARCH=arm64 go build -o dist/halptask-darwin-arm64 .
GOOS=darwin GOARCH=amd64 go build -o dist/halptask-darwin-amd64 .

# Linux (amd64 & ARM64)
GOOS=linux GOARCH=amd64 go build -o dist/halptask-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o dist/halptask-linux-arm64 .
```

---

## 📄 License

MIT License. Developed with Go & Charm Bubble Tea.

---

## 💡 Acknowledgments

Special thanks to the [LazyVim](https://github.com/LazyVim/LazyVim) project for inspiration.
