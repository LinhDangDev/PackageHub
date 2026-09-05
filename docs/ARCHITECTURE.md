# 🏗️ PackageHub Architecture & Engineering Design

## 1. System Overview

**PackageHub** is an ultra-fast, zero-bloat developer environment auditor, autonomous AI agent skills browser, and workstation synchronizer. It is built in **Go 1.22+** and packaged with an embedded **Microsoft Edge WebView2** native Windows GUI.

```
+-------------------------------------------------------------------------+
|                        PackageHub Desktop GUI                           |
|      (Cyber-Industrial Dark Theme • WebView2 Native Windows Window)     |
+------------------------------------+------------------------------------+
                                     | Local REST & Static Assets
+------------------------------------v------------------------------------+
|                         PackageHub Core (Go)                            |
|                                                                         |
|  +------------------------+  +-------------------+  +----------------+  |
|  |     Scanner Engine     |  |  Auditor & Health |  | Profile Engine |  |
|  | - Choco XML (.nuspec)  |  | - endoflife.date  |  | - YAML Export  |  |
|  | - Scoop App dirs       |  | - NPM Registry    |  | - ZIP Bundler  |  |
|  | - NPM package.json     |  | - PATH Resolver   |  | - Secret Mask  |  |
|  | - Agent Skills & Flags |  | - Engine Checker  |  | - Batch Install|  |
|  | - Project Dep Auditor  |  | - TTL Cache       |  | - Diff Engine  |  |
|  +------------------------+  +-------------------+  +----------------+  |
|                                                                         |
|  +-------------------------------------------------------------------+  |
|  |                        Installer & Auto-Fix                       |  |
|  | - Single Dependency Fixer (`npm install <pkg>@latest --save`)      |  |
|  | - Registry Environment Injector (`HKCU\Environment\Path`)         |  |
|  | - Version Switcher (Upgrade / Downgrade)                          |  |
|  | - Unattended Package Manager Execution (`CREATE_NO_WINDOW`)       |  |
|  +-------------------------------------------------------------------+  |
+-------------------------------------------------------------------------+
```

---

## 2. Core Subsystems

### 2.1. Hybrid Fast-Scanner (< 100ms)
Traditional developer environment scanners spawn sequential subprocesses (`node -v`, `git --version`, `npm list -g`). On Windows, `CreateProcess` introduces high kernel latency; scanning 50+ packages typically takes 15–25 seconds.

**PackageHub's Direct Filesystem Fast-Path:**
- **NPM Global**: Reads `%APPDATA%\npm\node_modules\*\package.json` directly.
- **Chocolatey**: Reads `C:\ProgramData\chocolatey\lib\*\*.nuspec` XML manifests directly.
- **Scoop**: Reads directory entries in `%USERPROFILE%\scoop\apps\*`.
- **AI Agent Skills**: Directly parses `SKILL.md` frontmatter and reads `.git/config` and `.git/HEAD` without executing `git.exe`.
- **Total Duration**: < 150 milliseconds for hundreds of packages and skills.

### 2.2. Zero-Console Window & DWM Immersive Dark Mode
- **Subsystem**: Compiled with `-ldflags="-H windowsgui"`, marking the PE binary as a native Windows GUI application (`IMAGE_SUBSYSTEM_WINDOWS_GUI`).
- **Zero Console Popups**: Windows does not spawn or allocate any console command prompt window when the app launches or runs background tasks.
- **AttachConsole**: When invoked from PowerShell or CMD with CLI flags (`scan`, `audit`), `kernel32.dll!AttachConsole` dynamically attaches output to the parent console.
- **Native Titlebar Synchronization**: Uses Desktop Window Manager (`dwmapi.dll`) to enable `DWMWA_USE_IMMERSIVE_DARK_MODE` (attribute 20) and `DWMWA_CAPTION_COLOR` (`#0b1019`), eliminating the stark white Windows titlebar.

### 2.3. Project & Drive Dependency Auditor
- Traverses codebases across specified drives (e.g. `D:\IdeaSideProject`) skipping noise directories (`node_modules`, `.git`, `.turbo`, `dist`).
- Identifies projects and language ecosystems:
  - **TypeScript / JavaScript**: `package.json`, `tsconfig.json`
  - **Go**: `go.mod`
  - **Rust**: `Cargo.toml`
  - **Python**: `requirements.txt`, `pyproject.toml`
- Audits dependencies against known deprecated packages (`request`, `moment`, `uuid@3`) and wildcard version risks (`latest`, `*`).
- Calculates health scores (0–100%) and provides **Granular Per-Dependency Auto-Fixing**.

### 2.4. Profile Engine & Offline ZIP Bundler
- **Declarative YAML (`packetinstall.yaml`)**:
  - Automatically scrubs and masks sensitive API credentials (`sk-ant-...`, `ghp_...`, `_TOKEN`, `_KEY`) into `${ENV_VAR}` placeholders.
- **Portable Offline ZIP Bundle**:
  - Compresses the YAML manifest + all local skill markdown files (`skills/*/SKILL.md`) + an automated `install.ps1` script into a single archive.
  - Can be transferred via USB to completely offline machines to restore workstation environments without git or external repositories.
