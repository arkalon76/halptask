# Agent Wiki Maintenance Guidelines 🤖

> **CRITICAL DIRECTIVE FOR FUTURE AI AGENTS**:
> Whenever you modify, add, or remove features, keybindings, CLI flags, configuration options, storage schemas, or UI workflows in **HalpTask**, you **MUST** update the relevant GitHub Wiki documentation in `wiki/`.

---

## 📌 Scope & Responsibilities

As an AI coding agent pair-programming on HalpTask, documentation is a primary codebase artifact. Code modifications are incomplete without updated documentation.

### When Code Changes, Update Wiki Pages:

1. **New or Modified Keybindings (`ui/keys.go`, `ui/app.go`)**:
   - Update `wiki/Keybindings-Reference.md` (tables for Leader Key `<space>` and Direct Vim motions).
   - Update `wiki/Power-User-Guide.md` if the change introduces a new motion or speed hack.
   - Update `docs/cheatsheet.md` and `README.md`.

2. **New Features or UI Components (e.g. Modals, Panels, Menus)**:
   - Update `wiki/Home.md` feature list and screenshot references.
   - Update `wiki/Getting-Started.md` quickstart tutorial if core workflows are impacted.
   - Update `wiki/User-Personas-&-Workflows.md` if the feature benefits specific roles (e.g. developers, DevOps, lawyers).

3. **Configuration Options (`config/config.go`, `ui/configmodal.go`)**:
   - Update `wiki/Configuration-&-Encryption-Security.md` (YAML key reference table & default values).

4. **Storage Format, Encryption, or CLI Flags (`model/storage.go`, `main.go`)**:
   - Update `wiki/Configuration-&-Encryption-Security.md` and `wiki/Getting-Started.md`.

---

## 📸 Screenshots & Visual Media Standards

- **Screenshots tell a better picture**: When adding or updating major UI components, ensure screenshot files in `wiki/images/` and `docs/images/` accurately depict current UI styling.
- Use `generate_image` or update snapshot tests when UI layouts evolve significantly.
- Always use standard GitHub Markdown image references: `![Alt Text](images/filename.jpg)` or raw repository links.

---

## 🔄 Wiki Synchronization Workflow

All Wiki files are maintained in the repository under `wiki/`:

```
wiki/
├── Home.md
├── Getting-Started.md
├── User-Personas-&-Workflows.md
├── Power-User-Guide.md
├── Configuration-&-Encryption-Security.md
├── Keybindings-Reference.md
├── Agent-Wiki-Maintenance-Guidelines.md
├── _Sidebar.md
├── _Footer.md
└── images/
    ├── halptask_banner.jpg
    ├── halptask_main.jpg
    ├── halptask_whichkey.jpg
    └── halptask_config.jpg
```

### Verification Checklist Before Ending Turn:
- [ ] Did I add/update keybinding tables in `Keybindings-Reference.md`?
- [ ] Is `AGENT.md` updated with references to new feature behavior?
- [ ] Do links between Wiki pages validate cleanly?
- [ ] Are screenshots and image references intact in `wiki/images/`?
