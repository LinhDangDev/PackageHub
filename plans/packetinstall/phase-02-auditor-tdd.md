---
phase: 2
title: "Auditor Core (TDD)"
status: pending
priority: P1
effort: "2h"
dependencies: [1]
---

# Phase 2: Auditor Core (TDD)

## Overview
Implement the obsolescence detection and upstream update auditing engine. Integrates the `endoflife.date` v1 REST API to identify deprecated/EOL runtimes (Node.js, Python, Go) and upstream package registry endpoints (NPM, PyPI, Crates.io) to detect available version upgrades.

## Requirements
- Query `endoflife.date/api/<product>.json` with cached HTTP client.
- Determine if a runtime major cycle has passed its EOL date.
- Query NPM registry `registry.npmjs.org/<pkg>/latest` for global CLI tools.
- Classify package health into:
  - `HEALTHY`: Current == Latest.
  - `OUTDATED`: Current < Latest (patch/minor/major).
  - `EOL_CRITICAL`: Runtime is past End-of-Life date.
- Local in-memory/file caching with TTL to prevent rate limits.

## TDD Test Cases & Fixtures
1. `TestEolChecker_Node18_IsEol`:
   - Mock HTTP response for Node.js EOL cycles with `httptest.Server`.
   - Node 18.20.8 (EOL date 2025-04-30) must be flagged as `EOL_CRITICAL`.
2. `TestEolChecker_Node22_IsActive`:
   - Node 22.14.0 (EOL date 2027-04-30) must be flagged as active.
3. `TestNpmUpdateChecker`:
   - Mock NPM registry response with latest version `2.1.261`.
   - Installed version `2.1.246` must report update available.
4. `TestAuditorCache`:
   - Ensure consecutive calls return cached responses without hitting HTTP backend.

## Success Criteria
- [ ] `go test -v ./internal/auditor` passes 100%.
- [ ] Correctly identifies EOL dates and generates actionable update commands.
