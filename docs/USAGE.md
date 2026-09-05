# User Guide

## Launching

### Desktop GUI

Double-click `packetinstall.exe` or run:

```powershell
.\packetinstall.exe
```

No console window appears — the app opens as a native desktop window with dark mode.

### CLI Mode

```powershell
.\packetinstall.exe scan             # Scan installed tools & agent skills
.\packetinstall.exe audit            # Check for EOL runtimes & upstream updates
.\packetinstall.exe export -o dev.yaml   # Export workstation profile to YAML
.\packetinstall.exe apply dev.yaml --dry-run  # Dry-run import from another machine
```

---

## Interface Guide

### Dashboard
High-level overview: environment health score, counts of packages/projects/skills/MCP servers, and AI coding agent status.

### System Care & Developer Optimizer (Tab: `System Care`)
A dedicated suite for deep maintenance:
1. **Zombie PATH Purger**:
   - Lists all entries in your Windows PATH with real-time existence validation.
   - Highlights dead, deleted, or inaccessible folders with red indicators.
   - Click **Prune Dead Paths** to cleanly remove invalid paths directly from the Windows User Registry.
2. **Dev Storage & Cache Hog Cleaner**:
   - Calculates disk usage across NPM cache, Python Pip cache, Go build/module caches, Rust Cargo, Yarn, pnpm, and Chocolatey.
   - Click **Clean** on any individual cache or **Free Up All Caches** to reclaim gigabytes of storage.
3. **Port & Process Conflict Auditor**:
   - Displays all active TCP listening developer ports (3000, 5173, 8080, 5432, etc.) along with their PID and process name (e.g. `node.exe`).
   - Click **Kill** to immediately free up a locked port without needing command-line `netstat` and `taskkill`.

### Dev Tools (Tab: `Dev Tools & Manage`)
Lists all installed packages (Chocolatey, Scoop, NPM Global) with:
- **Fix** — Auto-updates outdated packages or injects missing PATH entries.
- **Version ⇅** — Upgrade to latest/LTS or downgrade to a specific version.
- **Remove 🗑️** — Standard package manager uninstall.
- **Deep 💥** — Geek-style deep uninstall: executes package removal, scans for residual configuration files, registry keys, and dead PATH entries, and provides an interactive checklist to purge them permanently.

### Project Auditor
1. Enter a directory path in the scan input field (or use the quick-select drive buttons).
2. Click **Scan Projects Now**.
3. View discovered projects with health grades (A+ through D), language tags, and dependency counts.
4. Click a project card to inspect its dependencies. Use **Fix Dependency** to resolve a specific problematic package without touching others.

### Agent Skills
Browses all AI agent skills discovered on your machine:
- Filter by origin (Claude, OpenCode, Codex, Community) and complexity tier (Beginner, Intermediate, Pro).
- Click any skill to view its documentation, command syntax, and supported flags.

### MCP Servers
Displays Model Context Protocol servers configured in your editor settings.

### Machine Sync & Export
- **Export**: Generate a declarative YAML manifest or a portable offline ZIP bundle with selective checkboxes. Secrets are automatically masked.
- **Import**: Paste another machine's YAML to diff and auto-install missing tools.
