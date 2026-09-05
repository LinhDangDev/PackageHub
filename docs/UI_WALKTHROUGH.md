# UI Walkthrough

A visual guide to every screen and interaction in PackageHub.

---

## Top Navigation Bar

- **Brand & Status** — App name, version badge, and green connectivity indicator.
- **Search (`Ctrl+K`)** — Real-time filter across tools, skills, and projects.
- **Scan Latency** — Live display of total scan duration.
- **+ Install Tool** — Opens the in-app package installer modal.
- **Export Options** — Opens the selective export / ZIP bundler modal.

---

## Sidebar

Fixed navigation with tabs:
- ⚡ **Dashboard**
- 🛡️ **System Care** (with dynamic badge showing dead PATH count or CLEAN status)
- 📁 **Project Auditor**
- 📦 **Dev Tools & Manage**
- 🧠 **Agent Skills**
- 🔌 **MCP Servers**
- 🔄 **Machine Sync & Diff**

A pinned **Host Machine** card at the bottom shows OS info and online status.

---

## Dashboard

- **Health Score Banner** — Aggregate environment health (e.g. 98/100).
- **Metric Cards** — Clickable cards showing counts for packages, projects, skills, and MCP servers. Clicking navigates to the corresponding tab.
- **AI Agent Status** — Version badges and health pings for detected coding agents.

---

## System Care & Developer Optimizer (New!)

Provides three high-impact developer maintenance modules:

### 1. Zombie PATH Purger
- Scans all PATH entries and checks folder existence via `os.Stat`.
- Dead or inaccessible paths are highlighted with red warning badges (`⚠️ Dead Folder`).
- **Prune Dead Paths**: 1-click removes all dead entries from the Windows User Registry (`HKCU\Environment\Path`).

### 2. Dev Storage & Cache Hog Cleaner
- Audits reclaimable disk space across:
  - NPM Cache (`%APPDATA%\npm-cache`)
  - Python Pip Cache (`%LOCALAPPDATA%\pip\cache`)
  - Go Build Cache (`%LOCALAPPDATA%\go-build`)
  - Go Module Cache (`%USERPROFILE%\go\pkg\mod\cache`)
  - Rust Cargo Cache (`%USERPROFILE%\.cargo\registry\cache`)
  - Yarn & pnpm Caches
  - Chocolatey Cache
- Individual **Clean** buttons per cache or **Free Up All Caches** to reclaim gigabytes at once.

### 3. Port & Process Conflict Auditor
- Queries active TCP listening ports in the background.
- Associates listening ports with their process names and PIDs (e.g. `:3000 node.exe PID 1420`).
- **Kill**: Force-terminates the conflicting process to free up the port instantly.

---

## Dev Tools & Geek-Style Deep Uninstall

An interactive table of all installed packages with columns: Tool, Manager, Version, Path, Actions.

**Actions per tool:**
- **Fix** — 1-click update or PATH injection.
- **Version ⇅** — Opens a version switcher (latest, LTS, or custom).
- **Remove 🗑️** — Standard package uninstall.
- **Deep 💥 (New!)** — Opens the **Geek-Style Deep Uninstall Modal**:
  1. Runs package uninstallation.
  2. Deep-scans residual files, roaming data (`%APPDATA%`, `%LOCALAPPDATA%`), user dotfolders, registry keys (`HKCU\Software`), and dangling PATHs.
  3. Displays an interactive checklist with sizes.
  4. User reviews and clicks **Purge Selected Leftovers Permanently**.

---

## Project Auditor

1. Enter a path in the input field or use a quick-select drive button.
2. Click **Scan Projects Now**.
3. Projects appear as cards showing language, framework, health score, and dependency count.
4. Click a card to open the **Dependency Inspector Drawer**:
   - Lists every dependency with its declared version.
   - Flags deprecated packages and wildcard versions.
   - **Fix Dependency** resolves a single problematic package.
   - **Auto-Fix All** resolves every flagged issue in the project.
5. Click outside the drawer or press `ESC` to close.

---

## Agent Skills

Displays all discovered AI agent skills with:
- **Origin filters** — Claude, OpenCode, Codex, Community.
- **Tier filters** — Beginner, Intermediate, Pro.
- **Skill Cards** — Name, origin, tier, and short description. Click to open the documentation drawer.

---

## Machine Sync & Bundler

- **Current Manifest** — YAML view of your workstation config with secrets masked.
- **Import & Auto-Install** — Paste a YAML profile, view the diff, then run 1-click auto-install.
- **Export Options**:
  - **YAML** — Lightweight text manifest.
  - **ZIP Bundle** — Offline-portable archive with manifest, skill files, and install script.
  - Selective checkboxes let you pick exactly which tools/skills to include.
