---
phase: 3
title: "Profile Core (TDD)"
status: pending
priority: P1
effort: "2h"
dependencies: [1, 2]
---

# Phase 3: Profile Core (TDD)

## Overview
Implement the declarative profile engine (`packetinstall.yaml`) allowing developers to export their current machine state (runtimes, package manager packages, AI skills, and MCP configurations) into a portable manifest, and import/replay it on another machine with secret masking and dry-run diffs.

## Requirements
- Serialize scanned `SystemState` into structured YAML format.
- Mask sensitive environment variables in MCP server definitions (e.g. replace `sk-ant-xxx` with `${ANTHROPIC_API_KEY}`).
- Profile Importer:
  - Compare target machine state against imported profile.
  - Calculate `Diff`: Missing packages, Missing skills, Missing MCP servers.
  - Generate silent execution plan commands for Chocolatey, Scoop, Winget, and NPM.
  - Provide `--dry-run` output mode.

## TDD Test Cases & Fixtures
1. `TestProfileExport_SecretMasking`:
   - Input MCP server with `ANTHROPIC_API_KEY: "sk-ant-12345678"`.
   - Ensure exported YAML contains `${ANTHROPIC_API_KEY}` and zero raw secret bytes.
2. `TestProfileImport_DiffCalculation`:
   - Target machine has `git` and `node`.
   - Profile specifies `git`, `node`, `ripgrep`, and `@anthropic-ai/claude-code`.
   - Diff output must indicate `git` and `node` as skipped, and `ripgrep` + `claude-code` as pending install.
3. `TestExecutionPlanGeneration`:
   - Generate batch installer script or slice of commands with silent flags (`--yes`, `--silent`).

## Success Criteria
- [ ] `go test -v ./internal/profile` passes 100%.
- [ ] Zero secret leakage in exported manifests.
