# PackageHub

> Open-source developer environment manager, AI agent skills browser, project auditor, deep leftover cleaner & offline workstation bundler.

Built with Go and native Windows WebView2. Zero terminal popups, dark mode, scan in under 100ms.

[![Go](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows-lightgrey.svg)]()

---

## Features

- **Dashboard** — Overview of installed packages, project health, AI agent status.
- **Geek-Style Deep Uninstall** — Remove packages and permanently scavenge residual files (`%APPDATA%`, `%LOCALAPPDATA%`, dotfolders), registry keys (`HKCU\Software`), and dead PATH entries.
- **System Care & Optimizer**:
  - **Zombie PATH Purger** — Detect and prune dead/missing folders from Windows User PATH.
  - **Dev Cache & Storage Hog Cleaner** — One-click reclaim gigabytes from NPM, Pip, Go build/module, and Cargo caches.
  - **Port & Process Conflict Auditor** — Audit active TCP listening dev ports (3000, 5173, 8080, etc.) with PID and process name; 1-click terminate conflicting processes.
- **Project Auditor** — Scan drives/directories, detect languages & frameworks, audit dependency health (deprecations, wildcards), 1-click fix per dependency.
- **Dev Tools** — Manage Chocolatey, Scoop, and NPM global packages. Upgrade, downgrade, fix, or deep-remove.
- **Agent Skills Browser** — Auto-discovers AI agent skills across installed agent harnesses. Categorizes by origin and complexity tier.
- **MCP Servers** — Catalogs Model Context Protocol servers configured in your editor.
- **Machine Sync & Bundler** — Export your setup as YAML or a portable offline ZIP bundle. Import on a new machine with 1-click auto-install.

---

## Installation

### Prerequisites

- **Windows 10/11** with [Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)
- **Go 1.22+** (for building from source)

### Build from Source

```bash
git clone https://github.com/your-username/PackageHub.git
cd PackageHub
go build -ldflags="-H windowsgui" -o packetinstall.exe ./cmd/packetinstall
```

### Run

**Desktop GUI** — Double-click `packetinstall.exe` or:

```powershell
.\packetinstall.exe
```

**CLI mode:**

```powershell
.\packetinstall.exe scan             # Scan installed tools & skills
.\packetinstall.exe audit            # Check for EOL runtimes & updates
.\packetinstall.exe export -o dev.yaml   # Export workstation profile
.\packetinstall.exe apply dev.yaml --dry-run  # Preview import
```

---

## Usage

| Tab | What it does |
|-----|-------------|
| **Dashboard** | System health score, package/project/skill counts, AI agent status |
| **System Care** | Prune dead PATH folders, clean multi-GB dev caches (NPM, Pip, Go), and resolve port conflicts |
| **Dev Tools** | Installed packages table with Version ⇅, Fix, Remove, and **Deep Uninstall** |
| **Project Auditor** | Enter a path → scan projects → view dependency health grades (A+ to D) → fix individual deps |
| **Agent Skills** | Filter by origin (Claude, OpenCode, Codex, Community) and tier (Beginner/Intermediate/Pro) |
| **MCP Servers** | View configured MCP server connections |
| **Machine Sync** | Export/import workstation profiles. Secrets are auto-masked. |

For a detailed visual walkthrough, see [`docs/UI_WALKTHROUGH.md`](docs/UI_WALKTHROUGH.md).

---

## Project Structure

```
PackageHub/
├── cmd/packetinstall/    # Entrypoint (CLI + native window)
├── internal/
│   ├── app/              # WebView2 window & DWM dark mode
│   ├── model/            # Domain types
│   ├── cleaner/          # Leftovers scavenger, zombie PATH purger, dev cache & port auditor
│   ├── scanner/          # Choco, Scoop, NPM, Skills, MCP, Project scanners
│   ├── installer/        # Install, uninstall, version switch, auto-fix
│   ├── auditor/          # EOL & upstream update checker
│   ├── profile/          # YAML export/import, ZIP bundler, secret masking
│   └── web/              # REST API server & embedded UI
├── docs/                 # Documentation
└── LICENSE
```

---

## Documentation

- [UI Walkthrough](docs/UI_WALKTHROUGH.md)
- [User Guide](docs/USAGE.md)
- [Architecture](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Contributing](docs/CONTRIBUTING.md)

---

## License

[MIT](LICENSE)
