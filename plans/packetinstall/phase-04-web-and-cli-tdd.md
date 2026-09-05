---
phase: 4
title: "Web UI & CLI Integration (TDD)"
status: pending
priority: P1
effort: "2h"
dependencies: [1, 2, 3]
---

# Phase 4: Web UI & CLI Integration (TDD)

## Overview
Implement the unified entry point for `packetinstall`:
1. **CLI Subcommands:** `scan`, `audit`, `export`, `apply`, and `ui`.
2. **REST API Server:** Exposes endpoints `/api/scan`, `/api/audit`, `/api/profile/export`, `/api/profile/import`.
3. **Embedded Web Dashboard:** Single-page dashboard built with modern HTML5/CSS/Vanilla JS (or lightweight bundle) embedded directly into the Go executable via `//go:embed`.

## Requirements
- Standard Go `net/http` router with zero external heavy frameworks.
- Handlers:
  - `GET /api/scan`: Runs scanner and returns JSON state.
  - `GET /api/audit`: Runs auditor and returns version health status.
  - `POST /api/profile/export`: Exports YAML profile.
  - `POST /api/profile/apply`: Executes dry-run or live apply.
- Embedded static file handler at `/` serving the dashboard.
- CLI flags & subcommands:
  - `packetinstall scan` -> Print terminal table of installed tools/skills/MCPs.
  - `packetinstall audit` -> Print table with outdated / EOL warnings.
  - `packetinstall export -o my-profile.yaml` -> Export manifest.
  - `packetinstall apply my-profile.yaml --dry-run` -> Show diff plan.
  - `packetinstall ui` -> Start server and open default browser.

## TDD Test Cases & Fixtures
1. `TestApiScanEndpoint`:
   - Send `GET /api/scan` to `httptest.Server`.
   - Assert `200 OK` and JSON contains packages, skills, and MCP servers.
2. `TestApiAuditEndpoint`:
   - Send `GET /api/audit`.
   - Assert `200 OK` and JSON contains health audit items.
3. `TestEmbeddedStaticFiles`:
   - Send `GET /`.
   - Assert `200 OK` and response body contains dashboard HTML.

## Success Criteria
- [ ] `go test -v ./...` passes across all packages.
- [ ] `go build -o packetinstall.exe ./cmd/packetinstall` completes in < 2 seconds and outputs single binary.
