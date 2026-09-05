# ⚡ PackageHub (`packetinstall`)

> **Open-Source Developer Workstation, Autonomous AI Agent Skills, Project Auditor & Offline ZIP Bundler**  
> An ultra-fast, zero-bloat developer environment manager written in Go with an embedded native Windows WebView2 desktop interface.

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20Native%20Desktop%20%7C%20CLI-lightgrey.svg)]()
[![GitHub Repo](https://img.shields.io/badge/github-LinhDangDev%2FPackageHub-cyan.svg)](https://github.com/LinhDangDev/PackageHub)

---

## 🌟 Highlights & Key Features

### 1. 🧠 Autonomous AI Agent Skills Hub & VividKit-Style Docs
- **Tool Origin Separation**: Cleanly partitions skills by origin: **ClaudeKit (`ck:`)**, **OpenCode / OMP**, **Codex**, or **Community / Universal**.
- **VividKit Skill Tiers**:
  - 🟢 **Beginner**: `/ck:ask`, `/ck:brainstorm`, `/ck:fix`, `/ck:docs`, `/ck:interview-docs`...
  - 🟡 **Intermediate**: `/ck:plan`, `/ck:cook`, `/ck:test`, `/ck:review-pr`, `/ck:git`...
  - 🔴 **Pro**: `/ck:bootstrap`, `/ck:scout`, `/ck:ship`, `/ck:watzup`, `/ck:team`, `/ck:ai-artist`...
- **Detailed Command Flags & Parameter Explanations**:
  - Automatically analyzes arguments and flags (`--fast`, `--mode`, `--provider`, `--skip`, `--tdd`, `--parallel`, etc.) with clear, contextual explanations.
  - **`+ Add Flag` Button**: 1-click appends the flag directly into the sample command and copies it to your clipboard.
- **Interactive VividKit Documentation Drawer**:
  - Click any skill to slide out detailed documentation.
  - Features full command references, arguments, interactive 1-click copy examples, pro tips, and raw `SKILL.md` content.
  - **Click-Outside to Close**: Smoothly collapses when clicking anywhere outside the drawer or pressing `ESC`.

### 2. 📁 Drive & Project Dependency Auditor
- **Drive-Wide Codebase Discovery**: Recursively scans any drive or directory (e.g. `D:\IdeaSideProject`, `D:\`, `C:\Users\Dev`).
- **Language & Framework Detection**: Identifies **TypeScript / JavaScript** (Next.js, React, Turbo, Express), **Go** (Gin, Fiber), **Rust** (Cargo), and **Python**.
- **Dependency Health & Deprecation Audit**: Flags abandoned packages (`request`, `moment`, `uuid@3`) and wildcard version risks (`latest`, `*`).
- **Granular 1-Click Fix (`⚡ Fix Dependency`)**:
  - Resolves and pins **only the specific dependency that has an issue** (e.g. pinning `obsidian: "latest"` ➔ `npm install obsidian@latest --save`) without modifying other libraries.
- Calculates an objective project health score (0–100%, A+ to D) with an **`⚡ Auto-Fix All Project Dependencies`** button.

### 3. ⚡ Package Management & Auto-Fix Engine
- **`⚡ Fix` Button**: 1-click updates outdated packages or **automatically injects missing directories into the Windows User PATH Registry** (`HKCU\Environment\Path`).
- **`Version ⇅`**: Easily upgrade to `latest`, `LTS`, or downgrade to any specific version number.
- **`🗑️ Remove`**: Cleanly uninstalls packages via NPM Global, Chocolatey, or Scoop.
- **In-App Tool Installer**: Search and install tools directly within the UI without opening an external terminal.

### 4. 📦 Selective Export & Portable Offline ZIP Bundler
- **Selective Checklist**: Use checkboxes to select exact tools, skills, and MCP connectors to export.
- **Two Export Formats**:
  - 📄 **Declarative YAML (`packetinstall.yaml`)**: Lightweight manifest (~5KB) with sensitive credentials safely masked (`sk-ant-...` ➔ `${ENV_VAR}`).
  - 📦 **Portable Offline ZIP Bundle (`packetinstall-bundle.zip`)**: Packages the manifest, all offline skill markdown files (`skills/*/SKILL.md`), and an automated `install.ps1` script for offline air-gapped bootstrapping.

### 5. 🖥️ Zero-Console Windows Native Desktop Window
- **Pure Windows GUI**: Compiled with `-ldflags="-H windowsgui"`. When launched or double-clicked, **no black console command prompt window ever appears or flashes**.
- **DWM Dark Mode Synchronization**: Synchronizes titlebar color with the application palette (`#0b1019`) via Desktop Window Manager (`dwmapi.dll`).
- **Fixed & Sticky Sidebar**: Navigation and Host Machine status are 100% pinned; only the main content viewport scrolls.

---

## 🚀 Quick Start

### 1. Launch Native Windows Desktop App
Double-click `packetinstall.exe` in File Explorer or run in PowerShell:
```powershell
.\packetinstall.exe
```
*Launches the native desktop window immediately with dark mode titlebar and zero console windows.*

### 2. Command Line (CLI) Mode
```powershell
.\packetinstall.exe scan           # Fast CLI scan of tools & skills
.\packetinstall.exe audit          # Check for EOL runtimes & updates
.\packetinstall.exe export -o dev.yaml # Export machine profile
.\packetinstall.exe apply dev.yaml --dry-run
```

---

## 📚 Documentation

- [🖥️ UI & Feature Walkthrough](docs/UI_WALKTHROUGH.md) — Comprehensive visual guide of every screen, drawer, and button.
- [📖 User Guide & Operations Manual](docs/USAGE.md) — Step-by-step operating instructions.
- [🏗️ Architecture & Engineering Design](docs/ARCHITECTURE.md) — Deep dive into the Go core, hybrid scan, and Win32 internals.
- [🔌 REST API Reference](docs/API.md) — Full specification of all 12 local REST endpoints.
- [🤝 Contributing Guide](docs/CONTRIBUTING.md) — Guidelines for developer contributions and TDD tests.

---

## 🛠️ Repository Structure

```
PackageHub/
├── cmd/packetinstall/main.go     # CLI & Native Window entrypoint (AttachConsole)
├── internal/
│   ├── app/app.go                # Native Windows WebView2 Desktop Window & DWM Dark Mode
│   ├── model/types.go            # Domain models (Package, Project, Skill, McpServer, Profile)
│   ├── scanner/
│   │   ├── choco.go              # Chocolatey XML parser
│   │   ├── scoop.go              # Scoop apps parser
│   │   ├── npm.go                # Global NPM package.json parser
│   │   ├── skills.go             # Skills scanner with tool origin, tiers & flag parser
│   │   ├── mcp.go                # Claude Desktop & Cursor MCP parser
│   │   └── project.go            # Drive & Codebase dependency auditor
│   ├── installer/
│   │   ├── installer.go          # Install, Uninstall, Switch Version, Batch executor
│   │   └── fixer.go              # 1-Click Auto-Fix engine (Tools, PATH, Single Dependency)
│   ├── auditor/                  # EOL & Upstream registry update auditor
│   ├── profile/
│   │   ├── profile.go            # YAML export/import with secret masking
│   │   └── bundle.go             # Selective export & portable offline ZIP bundler
│   └── web/                      # HTTP REST API & Cyber-Industrial UI
├── docs/                         # Full English documentation suite
└── plans/                        # Architectural research & TDD roadmap
```

---

## 📄 License

MIT License. Open-source and free for all developers and creators.
