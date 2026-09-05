# 🖥️ PackageHub UI & Feature Walkthrough

This document provides a comprehensive, visual, and functional walkthrough of **PackageHub**'s Cyber-Industrial interface, detailing what every screen, drawer, modal, and button does.

---

## 1. Top Navigation Bar & Global Actions

```
+--------------------------------------------------------------------------------------------------------------------+
| [P] packetinstall v1.2 [Windows Native]    [Search 200+ tools, skills... Ctrl+K]    ⚡ 35ms  [+ Install]  [Export] |
+--------------------------------------------------------------------------------------------------------------------+
```

- **Brand Logo & Status Halo (`P`)**:
  - Displays the PackageHub branding with a pulsing green indicator signifying active backend connectivity.
  - Badges: `v1.2` release version and `Windows Native` (confirming Microsoft Edge WebView2 runtime mode).
- **Global Search Bar (`Ctrl + K`)**:
  - Instantly filters tools, AI skills, and code projects in real time (< 1ms).
  - Pressing `ESC` clears search and blurs input.
- **Scan Latency Badge (`⚡ 35ms`)**:
  - Live metric displaying the total filesystem walk and audit time for your workstation.
- **`+ Install Tool` Button**:
  - Opens the **In-App Package Installer Modal** to install tools via NPM, Chocolatey, or Scoop without switching to an external terminal.
- **`Export Options` Button**:
  - Opens the **Selective Export & Portable ZIP Bundler Modal** to generate custom YAML manifests or standalone offline ZIP archives.

---

## 2. Left Navigation Sidebar

- **Navigation Menu**:
  - ⚡ **Dashboard**: High-level system overview, health grades, and AI coding agent status cards.
  - 📁 **Project Auditor**: Drive-wide codebase scanner and granular dependency health auditor.
  - 📦 **Dev Tools & Manage**: Installed package table with version switcher, uninstaller, and diagnostics.
  - 🧠 **Agent Skills**: Categorized repository of 222+ AI agent skills with VividKit-style documentation.
  - 🔌 **MCP Servers**: Model Context Protocol servers configured in Claude Desktop and Cursor.
  - 🔄 **Machine Sync & Diff**: Cross-machine profile import, diff calculation, and 1-click auto-installer.
- **Host Machine Status Card (Pinned & Sticky)**:
  - Anchored permanently at the bottom of the sidebar.
  - Shows operating system (`Windows 11 Workstation`), online status, and Go core runtime (`Zero-Install`).
  - Remains strictly visible regardless of page scroll.

---

## 3. Screen Walkthroughs

### 3.1. Dashboard View
- **Hero Status Banner**:
  - Displays an aggregate environment health score (e.g. `98/100 Optimal`) with ambient cyan glow.
- **4 Interactive Metric Cards**:
  - **Dev Packages**: Total count of managed runtimes, Chocolatey, and NPM global packages (`52`).
  - **Scanned Projects**: Count of discovered codebases across chosen drives (`18 in D:\IdeaSideProject`).
  - **Agent Skills**: Total registered skills across agent harnesses (`222`).
  - **MCP Connectors**: Active server bindings in Claude and Cursor configs.
- **Autonomous AI Coding Agents Status Widget**:
  - Live health pings and detected version badges for **Claude Code** (`@anthropic-ai/claude-code`), **OpenCode / OMP** (`opencode-ai`), and **Codex / Gemini CLI**.

---

### 3.2. Project & Drive Dependency Auditor (Tab: `Project Auditor`)

Designed to audit real-world codebases across drives (e.g. `D:\IdeaSideProject`):

```
+-------------------------------------------------------------------------------------------------+
| 📁 Drive & Project Dependency Auditor                         [ D:\IdeaSideProject ]  [Scan Now]|
+-------------------------------------------------------------------------------------------------+
|                                                                                                 |
| [ dating-backend ]    [ readbook-monorepo ]   [ obsidian-chat ]      [ packetinstall ]          |
| TypeScript • Express  TypeScript • Turborepo  TypeScript • Node.js   Go • WebView Desktop       |
| Score: 100% (A+)      Score: 100% (A+)        Score: 88% (B) ⚠️1 dep  Score: 100% (A+)           |
| Deps: 19              Deps: 44                Deps: 13               Deps: 5                    |
|                                                                                                 |
+-------------------------------------------------------------------------------------------------+
```

#### What It Does:
1. **Auto-Detects Language & Framework**:
   - Analyzes project files to identify: **TypeScript / JavaScript** (Next.js, React, Turbo, Express), **Go** (Gin, Fiber, Desktop WebView), **Rust** (Cargo), **Python**.
2. **Audits Dependencies for Deprecations & Risks**:
   - Flags abandoned or deprecated packages (e.g. `request`, `moment`, `uuid@3`).
   - Flags loose wildcard version constraints (`latest`, `*`).
3. **Calculates Project Health Score**:
   - Assigns a score from 0 to 100% (A+ for clean projects, B or D for projects with vulnerabilities).
4. **Project Dependency Inspector Drawer**:
   - Clicking any project slides open a detailed drawer from the right.
   - Lists every `dependency` and `devDependency` with its declared version.
   - **Granular 1-Click Fix Button (`⚡ Fix Dependency`)**:
     - Targets **only the specific package that has an issue** (e.g. resolves `obsidian: "latest"` by running `npm install obsidian@latest --save`).
     - Updates `package.json` for that specific package and raises the project health score to 100% without modifying any other library!
   - **Click-Outside to Close**: Clicking outside the drawer (or pressing `ESC`) smoothly dismisses it.

---

### 3.3. Dev Tools & Package Management (Tab: `Dev Tools & Manage`)

- **Interactive Table**:
  - Columns: `Tool`, `Manager` (`choco`, `npm`, `scoop`), `Installed Version`, `Path`, `Actions`.
  - Tools are tagged with specialized branded icons (Python, Git, Docker, Neovim, Claude, Gemini, Codex).
- **Actions Available on Every Tool**:
  - **`⚡ Fix`**: 1-click automatic remediation (upgrades outdated tools or injects missing paths into the Windows User Registry `HKCU\Environment\Path`).
  - **`Version ⇅`**: Opens the **Version Switcher Modal** to upgrade to `latest`, `LTS`, or downgrade to any specific version number (e.g. `2.1.246`).
  - **`🗑️ Remove`**: Clean uninstaller running `npm uninstall -g`, `choco uninstall -y`, or `scoop uninstall`.

---

### 3.4. AI Agent Skills Hub & VividKit-Style Docs (Tab: `Agent Skills`)

```
+-------------------------------------------------------------------------------------------------+
| 🧠 AI Agent Skills & Commands Hub                                 Total Skills: 222             |
| [All (222)] [ClaudeKit (85)] [OpenCode/OMP (40)] [Codex (15)] [Community (82)]                  |
| Tier: [All] [🟢 Beginner] [🟡 Intermediate] [🔴 Pro]                                            |
+-------------------------------------------------------------------------------------------------+
|                                                                                                 |
| [ /ck:cook ]             [ /ck:plan ]             [ /ck:ask ]            [ /ck:ai-artist ]      |
| ClaudeKit • 🟡 Interm.    ClaudeKit • 🟡 Interm.    ClaudeKit • 🟢 Begin.  ClaudeKit • 🔴 Pro     |
| Structured implementation Technical roadmap       Architecture advice    Nano Banana visual art |
| [View Docs ➔]            [View Docs ➔]            [View Docs ➔]          [View Docs ➔]          |
|                                                                                                 |
+-------------------------------------------------------------------------------------------------+
```

#### What It Does:
1. **Tool Origin Separation**:
   - Skills are cleanly partitioned by ecosystem: **ClaudeKit (`ck:`)**, **OpenCode / OMP**, **Codex**, and **Universal / Community**.
2. **VividKit Skill Tiers**:
   - Organizes ClaudeKit commands into:
     - 🟢 **Beginner**: `/ck:ask`, `/ck:brainstorm`, `/ck:fix`, `/ck:docs`
     - 🟡 **Intermediate**: `/ck:plan`, `/ck:cook`, `/ck:test`, `/ck:review-pr`, `/ck:git`
     - 🔴 **Pro**: `/ck:bootstrap`, `/ck:scout`, `/ck:ship`, `/ck:watzup`, `/ck:team`, `/ck:ai-artist`
3. **Interactive VividKit Documentation Drawer**:
   - Clicking any skill card slides open the full documentation drawer:
     - **Command Reference**: Full syntax with highlighted arguments (e.g. `/ck:cook [task] [--fast|--parallel|--tdd]`).
     - **Command Flags & Parameters**:
       - Lists every supported flag (`--fast`, `--mode`, `--provider`, `--skip`, `--tdd`, `--parallel`, `--advice`, etc.).
       - Displays clear, concise explanations for each flag's behavior.
       - **`+ Add Flag` Button**: 1-click appends the selected flag into your command prompt and copies it to clipboard!
     - **Interactive Example Prompts**: Real-world prompts with 1-click copy buttons.
     - **Pro Tips from ClaudeKit**: Expert recommendations for token management, context files (`@plan.md`), and memory clearing (`/clear`).
     - **Full Markdown Docs**: Complete unredacted `SKILL.md` content.
   - **Click-Outside to Close**: Click anywhere outside the drawer to collapse it.

---

### 3.5. Machine Sync & Offline ZIP Bundler (Tab: `Machine Sync & Diff`)

```
+-------------------------------------------------------------------------------------------------+
| Left: Current Machine YAML                         Right: Compare & Auto-Install                |
| [ Copy YAML ]                                      [ Paste packetinstall.yaml ]                 |
|                                                    [ 🚀 Execute Auto-Install All (5 items) ]    |
| schema_version: "1.0"                              -------------------------------------------- |
| runtimes: ...                                      [1/5] choco install -y ripgrep... ✓ Done     |
| global_clis: ...                                   [2/5] npm install -g @anthropic-ai... ✓ Done |
+-------------------------------------------------------------------------------------------------+
```

#### What It Does:
1. **Current Workstation Manifest**:
   - Renders the local machine configuration with sensitive secrets masked (`sk-ant-...` ➔ `${ENV_VAR}`).
2. **1-Click Auto-Install All**:
   - When importing a profile from another machine, diffs are calculated automatically.
   - Clicking **`🚀 Execute Auto-Install All`** unattendedly runs the installation queue in the background with real-time step logging and auto-refresh.
3. **Export Options (Top Navigation Bar)**:
   - **Declarative YAML (`packetinstall.yaml`)**: Lightweight manifest (~5KB).
   - **Portable Offline ZIP Bundle (`packetinstall-bundle.zip`)**: Packages the manifest, all local skill files (`skills/*/SKILL.md`), and an automated `install.ps1` script for offline, air-gapped machine bootstrap.
   - **Selective Checkboxes**: Choose exactly which tools and skills to include.
