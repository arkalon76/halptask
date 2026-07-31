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
- **[Release Guide](docs/RELEASE.md)**: Instructions for the automated GoReleaser pipeline.

---

## ✨ Features

- **Vim Native Keybindings**: Native Vim navigation (`j`, `k`, `h`, `l`, `gg`, `G`, `oo`, `oc`, `O`, `dd`, `x`, `u`, `ctrl+r`, `tab`, `shift+tab`).
- **Leader Menu (`<space>`)**: Leader key popup window displaying all available shortcuts non-intrusively.
- **Dynamic WhichKey Popup**: Visual popup updates as you type key prefixes for Leader options (`<space> b`, `<space> t`, `<space> z`, `<space> e`) as well as all multi-character Vim prefixes (`o`, `d`, `z`, `g`, `w`, `f`).
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
auto_save: true
data_file: ~/.config/halptask/data.txt
encrypted: false
indent_spaces: 2
leader_key: " "
show_which_key: true
theme: default
```

**Auto-Save & Encrypted Files**: 
When `auto_save: true`, HalpTask will automatically save all tree state mutations in the background. If you open or create an encrypted file but haven't provided a passphrase yet, auto-save will pause until you enter your passphrase to prevent data loss or lockouts.

---

## 🛠️ Releases & Compilation

We use [GoReleaser](https://goreleaser.com) and GitHub Actions to automatically build and release binaries for macOS, Linux, and Windows across `amd64` and `arm64` architectures.

For full details on the automated release process or how to build releases locally, refer to the **[Release Guide](docs/RELEASE.md)**.

---

## 📄 License

MIT License. Developed with Go & Charm Bubble Tea.

---

## 💡 Acknowledgments

Special thanks to the [LazyVim](https://github.com/LazyVim/LazyVim) project for inspiration.
