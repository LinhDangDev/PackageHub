# Architecture

## Overview

PackageHub is a Go application that scans, audits, and manages developer tools, AI agent skills, project dependencies, and system storage hygiene. It embeds a native Windows WebView2 GUI and a local REST API server.

```
┌─────────────────────────────────────────────────────┐
│                  WebView2 Desktop GUI                │
│         (Cyber-Industrial dark theme, DWM sync)      │
├──────────────────────┬──────────────────────────────┤
│   REST API Server    │   Embedded Static Assets     │
├──────────────────────┴──────────────────────────────┤
│                     Go Core                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │ Scanner  │  │ Auditor  │  │ Profile Engine   │   │
│  │ • Choco  │  │ • EOL    │  │ • YAML Export    │   │
│  │ • Scoop  │  │ • NPM    │  │ • ZIP Bundler    │   │
│  │ • NPM    │  │ • PATH   │  │ • Secret Masking │   │
│  │ • Skills │  │ • Cache  │  │ • Diff & Import  │   │
│  │ • MCP    │  │          │  │                  │   │
│  │ • Project│  │          │  │                  │   │
│  └──────────┘  └──────────┘  └──────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │          Cleaner & System Care               │   │
│  │  • Geek-Style Leftovers Scavenger            │   │
│  │  • Zombie PATH Purger                        │   │
│  │  • Dev Storage & Cache Cleaner               │   │
│  │  • Port & Process Conflict Auditor           │   │
│  └──────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │            Installer & Auto-Fix              │   │
│  │  • Package install/uninstall/version switch  │   │
│  │  • Single-dependency fixer                   │   │
│  │  • Windows PATH registry injection           │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Key Subsystems

### 1. Direct Filesystem Scanning (< 100ms)
Instead of spawning subprocesses (`node -v`, `git --version`), the scanner reads package manifests directly from disk:
- **Chocolatey**: `.nuspec` XML files in `%ProgramData%\chocolatey\lib\`
- **Scoop**: App directories in `%USERPROFILE%\scoop\apps\`
- **NPM Global**: `package.json` files in `%APPDATA%\npm\node_modules\`
- **Agent Skills**: `SKILL.md` frontmatter in standard skill directories

### 2. Cleaner & System Care Subsystem
- **Leftovers Scavenger**: When deep-uninstalling a package, scans for roaming/local data (`%APPDATA%`, `%LOCALAPPDATA%`), registry keys (`HKCU\Software`), and dead PATH entries.
- **Zombie PATH Purger**: Checks all PATH segments against `os.Stat(p)`. If invalid or deleted, strips them from the Windows User Registry (`HKCU\Environment\Path`) and broadcasts environment change.
- **Dev Cache Cleaner**: Calculates recursive disk usage across NPM, Pip, Go build (`%LOCALAPPDATA%\go-build`), Go module (`%USERPROFILE%\go\pkg\mod\cache`), and Cargo caches, offering granular or bulk cleaning.
- **Port Conflict Auditor**: Enumerates active listening TCP ports via background `netstat -ano -p tcp` and maps PIDs to process names via a single cached `tasklist` call.

### 3. Zero-Console Window & DWM Dark Mode
Compiled with `-ldflags="-H windowsgui"` (PE subsystem `WINDOWS_GUI`). No console window ever appears. When invoked from a terminal with CLI flags, `kernel32.dll!AttachConsole` attaches to the parent console for output. The titlebar is synchronized with the app theme via `dwmapi.dll`.

### 4. Secret Masking
The profile exporter detects sensitive values (API keys, tokens) and replaces them with `${ENV_VAR}` placeholders before writing YAML.

## Module Map

| Directory | Responsibility |
|-----------|---------------|
| `cmd/packetinstall/` | CLI parsing, AttachConsole, launches GUI or CLI mode |
| `internal/app/` | WebView2 window creation, DWM dark mode |
| `internal/model/` | Shared domain types (Package, Skill, Project, Profile, Cleaner) |
| `internal/cleaner/` | Leftovers scavenger, zombie PATH purger, dev cache & port auditor |
| `internal/scanner/` | Filesystem scanners for all package managers and skills |
| `internal/installer/` | Install, uninstall, version switch, auto-fix engine |
| `internal/auditor/` | EOL checks via endoflife.date, NPM registry lookups |
| `internal/profile/` | YAML/ZIP export, import, diff calculation, secret scrubbing |
| `internal/web/` | HTTP server, REST API handlers, embedded static UI |
