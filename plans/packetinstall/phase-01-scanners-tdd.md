---
phase: 1
title: "Scanner Core (TDD)"
status: pending
priority: P1
effort: "2h"
dependencies: []
---

# Phase 1: Scanner Core (TDD)

## Overview
Implement the core domain models and high-speed filesystem parsers for Chocolatey, Scoop, NPM Globals, AI Agent Skills, and Model Context Protocol (MCP) server configurations. All scanning operations must run concurrently via Goroutines and return within 100ms.

## Requirements
- Parse Chocolatey `.nuspec` XML files in `C:\ProgramData\chocolatey\lib\<pkg>\<pkg>.nuspec`.
- Parse Scoop apps directory structure and read symlink targets in `~/scoop/apps/<app>/current`.
- Parse global NPM packages (including scoped packages like `@anthropic-ai/claude-code`) from `package.json` in global node modules.
- Parse Agent Skills (`SKILL.md` frontmatter) in `~/.agent/skills/` and `~/.omp/agent/skills/`.
- Parse MCP server JSON configs for Claude Desktop (`claude_desktop_config.json`) and Cursor (`mcp.json`).

## Architecture & Data Flow
```
ScanCoordinator.ScanAll()
  ├── go scanChoco(baseDir)
  ├── go scanScoop(baseDir)
  ├── go scanNpm(baseDir)
  ├── go scanSkills(baseDir)
  └── go scanMcp(claudeConfig, cursorConfig)
  └── Collect into Unified SystemState struct
```

## TDD Test Cases & Fixtures
1. `TestChocoScanner`:
   - Setup temporary directory with mock `ripgrep.nuspec`.
   - Verify `ID == "ripgrep"` and `Version == "14.1.0"`.
2. `TestScoopScanner`:
   - Setup temporary directory with mock `~/scoop/apps/fzf/1.2.3`.
   - Verify detected app `fzf` with version `1.2.3`.
3. `TestNpmGlobalScanner`:
   - Setup mock `node_modules` with standard and scoped (`@anthropic-ai/claude-code`) packages.
   - Verify parsed packages and versions.
4. `TestSkillsScanner`:
   - Setup mock `skills/sequential-thinking/SKILL.md`.
   - Verify skill name, description, and git remote resolution.
5. `TestMcpScanner`:
   - Setup mock `claude_desktop_config.json`.
   - Verify server name, command, args, and env vars.

## Success Criteria
- [ ] `go test -v ./internal/scanner` passes 100%.
- [ ] Direct filesystem scan of local host executes in < 100ms.
