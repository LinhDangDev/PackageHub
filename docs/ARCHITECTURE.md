# Architecture

## Overview

PackageHub is a Go application that scans, audits, and manages developer tools, AI agent skills, and project dependencies. It embeds a native Windows WebView2 GUI and a local REST API server.

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
│  │            Installer & Auto-Fix              │   │
│  │  • Package install/uninstall/version switch  │   │
│  │  • Single-dependency fixer                   │   │
│  │  • Windows PATH registry injection           │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Key Design Decisions

### Direct Filesystem Scanning (< 100ms)
Instead of spawning subprocesses (`node -v`, `git --version`), the scanner reads package manifests directly from disk:
- **Chocolatey**: `.nuspec` XML files in `%ProgramData%\chocolatey\lib\`
- **Scoop**: App directories in `%USERPROFILE%\scoop\apps\`
- **NPM Global**: `package.json` files in `%APPDATA%\npm\node_modules\`
- **Agent Skills**: `SKILL.md` frontmatter in standard skill directories

This avoids `CreateProcess` overhead and completes the full scan in under 100ms.

### Zero-Console Window
The binary is compiled with `-ldflags="-H windowsgui"` (PE subsystem `WINDOWS_GUI`). No console window ever appears. When invoked from a terminal with CLI flags, `kernel32.dll!AttachConsole` attaches to the parent console for output.

### DWM Dark Mode
The titlebar is synchronized with the app theme via `dwmapi.dll` (`DWMWA_USE_IMMERSIVE_DARK_MODE` and `DWMWA_CAPTION_COLOR`).

### Secret Masking
The profile exporter detects sensitive values (API keys, tokens) and replaces them with `${ENV_VAR}` placeholders before writing YAML.

## Module Map

| Directory | Responsibility |
|-----------|---------------|
| `cmd/packetinstall/` | CLI parsing, AttachConsole, launches GUI or CLI mode |
| `internal/app/` | WebView2 window creation, DWM dark mode |
| `internal/model/` | Shared domain types (Package, Skill, Project, Profile) |
| `internal/scanner/` | Filesystem scanners for all package managers and skills |
| `internal/installer/` | Install, uninstall, version switch, auto-fix engine |
| `internal/auditor/` | EOL checks via endoflife.date, NPM registry lookups |
| `internal/profile/` | YAML/ZIP export, import, diff calculation, secret scrubbing |
| `internal/web/` | HTTP server, REST API handlers, embedded static UI |
