# REST API Reference

PackageHub runs a local HTTP server. All endpoints are served on `http://127.0.0.1:<port>`.

## Endpoints

### Scanning & Discovery

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/scan` | Scan all installed packages (Chocolatey, Scoop, NPM Global) |
| `GET` | `/api/skills` | Discover all AI agent skills |
| `GET` | `/api/mcp` | List configured MCP servers |
| `POST` | `/api/projects/scan` | Scan a directory for code projects |

**POST `/api/projects/scan`** body:
```json
{
  "path": "D:/Projects",
  "depth": 3
}
```

### Package Management

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/packages/install` | Install a package |
| `POST` | `/api/packages/uninstall` | Uninstall a package |
| `POST` | `/api/packages/switch-version` | Change package version |
| `POST` | `/api/fix` | Auto-fix a package issue |
| `POST` | `/api/projects/fix` | Fix a single project dependency |

### Profile & Export

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/profile/export` | Export workstation profile as YAML |
| `POST` | `/api/profile/export-selective` | Export selected items only |
| `POST` | `/api/profile/export-zip` | Export as portable offline ZIP bundle |
| `POST` | `/api/profile/diff` | Calculate diff between current state and a profile |
| `POST` | `/api/profile/export-batch` | Batch install from a profile |
