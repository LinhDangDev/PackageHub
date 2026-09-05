# 🔌 PackageHub REST API Specification

All endpoints are hosted locally by the embedded Go HTTP server on `127.0.0.1:<allocated-port>` (default CLI port: `3456`).

---

## 1. System & Discovery

### `GET /api/scan`
Runs the fast-path hybrid scanner across Chocolatey, Scoop, NPM Globals, AI Skills, and MCP configs.

**Response `200 OK`:**
```json
{
  "packages": [
    {
      "manager": "npm",
      "name": "@anthropic-ai/claude-code",
      "version": "2.1.246",
      "path": "C:\\Users\\Dev\\AppData\\Roaming\\npm\\node_modules\\@anthropic-ai\\claude-code"
    }
  ],
  "skills": [
    {
      "name": "ck:cook",
      "tool_origin": "claudekit",
      "tier": "intermediate",
      "command": "/ck:cook",
      "argument_hint": "[task] [--fast|--parallel|--tdd]",
      "description": "Implement features with structured workflow."
    }
  ],
  "mcp_servers": [],
  "scanned_at": "2026-09-05T21:40:00Z",
  "duration_ms": 38
}
```

---

## 2. Project & Dependency Auditor

### `POST /api/projects/scan`
Scans a target directory or drive for code projects and dependency health.

**Request Body:**
```json
{
  "path": "D:/IdeaSideProject",
  "depth": 3
}
```

**Response `200 OK`:**
```json
{
  "scan_path": "D:/IdeaSideProject",
  "projects": [
    {
      "name": "my-web-app",
      "path": "D:/IdeaSideProject/my-web-app",
      "language": "TypeScript",
      "framework": "Next.js",
      "health_score": 88,
      "dependencies": [
        {
          "name": "obsidian",
          "version": "latest",
          "status": "RISKY",
          "issue": "Wildcard version constraint."
        }
      ],
      "outdated_count": 0,
      "issue_count": 1
    }
  ],
  "total_deps": 237,
  "duration_ms": 54
}
```

### `POST /api/project/fix-dep`
Fixes and resolves a single specific dependency in a project without modifying other packages.

**Request Body:**
```json
{
  "project_path": "D:/IdeaSideProject/my-web-app",
  "language": "TypeScript",
  "package_name": "obsidian",
  "is_dev": false
}
```

---

## 3. Package Management & Auto-Fix

### `POST /api/packages/install`
Installs a new package via NPM, Chocolatey, or Scoop.

**Request Body:**
```json
{
  "manager": "npm",
  "package": "tsx",
  "global": true
}
```

### `POST /api/packages/uninstall`
Uninstalls a package.

**Request Body:**
```json
{
  "manager": "npm",
  "package": "tsx"
}
```

### `POST /api/packages/switch-version`
Upgrades or downgrades a package to a specific version.

**Request Body:**
```json
{
  "manager": "npm",
  "package": "@anthropic-ai/claude-code",
  "version": "2.1.250",
  "global": true
}
```

### `POST /api/fix`
Executes automated remediation for a tool, missing PATH, or project.

**Request Body:**
```json
{
  "type": "tool",
  "manager": "npm",
  "package": "@anthropic-ai/claude-code"
}
```

---

## 4. Profile Sync & Offline Bundler

### `POST /api/profile/diff`
Compares an imported YAML profile against the target machine.

### `POST /api/profile/execute-batch`
Executes an array of batch installation commands sequentially.

### `POST /api/profile/export-zip`
Generates and downloads a portable offline ZIP bundle containing `packetinstall.yaml`, offline skill files, and `install.ps1`.

### `POST /api/profile/export-selective`
Generates and downloads a custom YAML profile containing only user-selected items.
