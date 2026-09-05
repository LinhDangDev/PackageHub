# 📖 PackageHub User Guide & Operations Manual

For an in-depth visual breakdown of every screen, button, and interaction, see the [Full UI Walkthrough](UI_WALKTHROUGH.md).

---

## 1. Launching PackageHub

### 1.1. Native Desktop GUI (Default)
Double-click `packetinstall.exe` in File Explorer or run in PowerShell:
```powershell
.\packetinstall.exe
```
*Opens the standalone desktop application window with Dark Mode titlebar and zero console popups.*

### 1.2. Command Line (CLI) Subcommands
```powershell
# Scan installed packages & skills in terminal
.\packetinstall.exe scan

# Audit for End-of-Life runtimes & upstream updates
.\packetinstall.exe audit

# Export workstation configuration
.\packetinstall.exe export -o workstation.yaml

# Compare and dry-run apply a profile from another machine
.\packetinstall.exe apply workstation.yaml --dry-run
```

---

## 2. Navigating the Interface

### 2.1. Dashboard View
- **Hero Status Banner**: Environment health assessment (e.g. `98/100 Optimal`).
- **Metric Cards**: Real-time counts of installed packages, discovered projects, agent skills, and active MCP servers.
- **Autonomous AI Coding Agents Widget**: Version status and health pings for Claude Code, OpenCode / OMP, and Codex / Gemini CLI.

### 2.2. Project & Drive Dependency Auditor (Tab: `Project Auditor`)
1. Enter or select any directory/drive (e.g. `D:\IdeaSideProject`, `D:\`, `C:\Users\Dev`).
2. Click **"Scan Projects Now"**.
3. View discovered projects with health grades (`A+ 100%`, `B`, `D`), language tags, and issue counts.
4. Click any project card to open the **Project Dependency Inspector Drawer**:
   - Inspect every declared package and version.
   - For risky wildcard versions (`latest`, `*`) or deprecated packages, click the dedicated **"⚡ Fix Dependency"** button to resolve and pin that specific package without affecting other dependencies.
   - Click outside the drawer or press `ESC` to close.

### 2.3. Dev Tools & Package Management (Tab: `Dev Tools & Manage`)
- **Upgrade / Downgrade Version Switcher**: Click the **Version ⇅** button on any package to upgrade to `latest`, `LTS`, or enter any specific version number.
- **Uninstalling Tools**: Click the trash icon 🗑️ to remove a tool cleanly from NPM, Chocolatey, or Scoop.
- **Installing New Tools**: Click **"+ Install Tool"** in the top navigation bar to search and install tools with live console output.

### 2.4. AI Agent Skills Hub & VividKit-Style Docs (Tab: `Agent Skills`)
- **Ecosystem Origin Filters**: Switch between `ClaudeKit`, `OpenCode / OMP`, `Codex`, and `Community / Universal`.
- **VividKit Skill Tiers**:
  - 🟢 **Beginner**: `/ck:ask`, `/ck:brainstorm`, `/ck:fix`, `/ck:docs`
  - 🟡 **Intermediate**: `/ck:plan`, `/ck:cook`, `/ck:test`, `/ck:review-pr`, `/ck:git`
  - 🔴 **Pro**: `/ck:bootstrap`, `/ck:scout`, `/ck:ship`, `/ck:watzup`, `/ck:team`, `/ck:ai-artist`
- **Interactive Documentation Drawer**:
  - Click any skill to view its command syntax, arguments, interactive 1-click copy examples, pro tips, and full Markdown documentation.
  - **Flag Explanations & 1-Click Append**: Inspect every supported flag (`--fast`, `--mode`, `--provider`, `--tdd`, `--skip`) with detailed descriptions. Click `+ Add Flag` to append the flag directly to your prompt.
  - Click outside the drawer or press `ESC` to dismiss.

### 2.5. Machine Sync & Offline ZIP Bundler (Tab: `Machine Sync & Diff`)
- **1-Click Auto-Install All**: Paste another machine's `packetinstall.yaml` to calculate differences. Click **"🚀 Execute Auto-Install All"** to execute unattended background installations with live progress tracking.
- **Export Options**:
  - Click **"Export Options"** in the top navigation bar.
  - Choose between **Declarative YAML** (lightweight text manifest) or **Portable Offline ZIP Bundle** (packages manifest + offline skill files + `install.ps1` script).
  - Use checkboxes to select only specific packages or skills to include.
