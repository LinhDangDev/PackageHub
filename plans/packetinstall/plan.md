# Master Plan: packetinstall (Open-Source Dev Ecosystem & Agent Management)

**Architecture:** Go 1.27 Core Engine + Embedded Web UI / CLI  
**Methodology:** Strict Test-Driven Development (TDD)  
**Target Platform:** Windows First (Chocolatey, Scoop, Winget, NPM Global, Pipx, Cargo, MCP, Agent Skills) & Cross-Platform (macOS, Linux)

---

## 1. Executive Overview

`packetinstall` is a high-speed, lightweight developer environment auditor and profile synchronizer inspired by the STM (Smart Tools Management) workflow. It provides:
1. **Hybrid Fast Scanning (< 100ms):** Direct filesystem inspection of Chocolatey (`.nuspec`), Scoop (`apps/`), NPM Global (`package.json`), Cargo (`.crates2.json`), Pipx (`pipx_metadata.json`), AI Agent Skills (`~/.agent/skills/`), and MCP Server configs.
2. **Obsolescence & Update Auditing:** Integrates `endoflife.date` API to flag End-of-Life (EOL) runtimes (Node, Python, Go) and upstream package registries (NPM, PyPI, Crates.io) for update availability.
3. **Declarative Profile Engine:** Export system state into `packetinstall.yaml` with secret scrubbing, and import across machines with dry-run verification.
4. **Dual Interface:** Full headless CLI (`packetinstall scan/audit/export/apply`) plus an embedded, zero-install Web UI (`packetinstall ui`) served via Go standard `net/http` with embedded assets.

---

## 2. Test-Driven Development (TDD) Strategy

Every feature follows the strict Red-Green-Refactor cycle:
1. **Red:** Write unit tests in `*_test.go` with mock directory fixtures and HTTP test servers (`httptest.NewServer`) asserting expected outputs. Run `go test` -> Fails.
2. **Green:** Implement the minimal production code in `internal/*`. Run `go test` -> Passes.
3. **Refactor:** Optimize performance, clean up error handling, and ensure zero allocations where possible.

---

## 3. Project Directory Structure

```
packetinstall/
├── cmd/
│   └── packetinstall/
│       └── main.go                 # CLI commands & Server launcher
├── internal/
│   ├── model/
│   │   └── types.go                # Core domain structs (Package, Skill, McpServer, Profile)
│   ├── scanner/
│   │   ├── scanner.go              # Coordinator (parallel scanning via goroutines)
│   │   ├── choco.go                # Choco nuspec XML parser
│   │   ├── choco_test.go           # Choco unit tests
│   │   ├── scoop.go                # Scoop apps directory parser
│   │   ├── scoop_test.go           # Scoop unit tests
│   │   ├── npm.go                  # NPM global package.json parser
│   │   ├── npm_test.go             # NPM unit tests
│   │   ├── skills.go               # AI Agent Skills (SKILL.md & git remote) parser
│   │   ├── skills_test.go          # Skills unit tests
│   │   ├── mcp.go                  # Claude Desktop & Cursor MCP config parser
│   │   └── mcp_test.go             # MCP unit tests
│   ├── auditor/
│   │   ├── auditor.go              # Coordinator for version & EOL checks
│   │   ├── eol.go                  # endoflife.date REST client
│   │   ├── eol_test.go             # EOL client unit tests (httptest)
│   │   ├── registry.go             # NPM / PyPI / Crates latest version checks
│   │   ├── registry_test.go        # Registry check unit tests
│   │   └── cache.go                # Memory / file cache with TTL
│   ├── profile/
│   │   ├── profile.go              # packetinstall.yaml schema definition
│   │   ├── exporter.go             # State serialization with secret scrubbing
│   │   ├── exporter_test.go        # Exporter unit tests
│   │   ├── importer.go             # Diff calculation, dry-run & command plan
│   │   └── importer_test.go        # Importer unit tests
│   └── web/
│       ├── server.go               # HTTP REST API server
│       ├── server_test.go          # HTTP API endpoint tests
│       └── static/                 # Embedded HTML/CSS/JS single-page dashboard
├── plans/
│   └── packetinstall/
│       ├── plan.md
│       ├── phase-01-scanners-tdd.md
│       ├── phase-02-auditor-tdd.md
│       ├── phase-03-profile-tdd.md
│       └── phase-04-web-and-cli-tdd.md
├── go.mod
└── go.sum
```

---

## 4. Implementation Phases

| Phase | Focus | Key Deliverables | Verification |
|---|---|---|---|
| **Phase 1: Scanner Core** | Fast-path filesystem scanning | Models, Choco, Scoop, NPM, Skills, MCP parsers | `go test -v ./internal/scanner` |
| **Phase 2: Auditor Core** | Update & EOL detection | EOL client, Registry client, local cache | `go test -v ./internal/auditor` |
| **Phase 3: Profile Core** | Declarative manifest sync | YAML export (secret scrubbing), Import diff & install runner | `go test -v ./internal/profile` |
| **Phase 4: Web UI & CLI** | User interfaces & Integration | Embedded Web dashboard, REST API, CLI subcommands | `go test -v ./...` + `go build` |
