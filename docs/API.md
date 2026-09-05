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

### Package Management & Deep Uninstall

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/packages/install` | Install a package |
| `POST` | `/api/packages/uninstall` | Standard uninstall of a package |
| `POST` | `/api/packages/switch-version` | Change package version |
| `POST` | `/api/fix` | Auto-fix a package issue |
| `POST` | `/api/projects/fix` | Fix a single project dependency |

### System Care & Cleaner

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/cleaner/leftovers?tool=<name>` | Scan residual files, registry keys, and paths for a tool |
| `POST` | `/api/cleaner/purge-leftovers` | Permanently purge selected leftovers |
| `GET` | `/api/cleaner/path` | Audit Windows PATH and detect dead directories |
| `POST` | `/api/cleaner/path-prune` | Prune all dead paths from Windows User Registry |
| `GET` | `/api/cleaner/caches` | Audit disk usage across developer caches (NPM, Pip, Go, Cargo) |
| `POST` | `/api/cleaner/cache-clean` | Clean a specific developer cache (`{"cache_id": "npm"}`) |
| `GET` | `/api/cleaner/ports` | List active TCP listening dev ports with PIDs and process names |
| `POST` | `/api/cleaner/kill-port` | Force-terminate a process by PID (`{"pid": 1420}`) |

### Profile & Export

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/profile/export` | Export workstation profile as YAML |
| `POST` | `/api/profile/export-selective` | Export selected items only |
| `POST` | `/api/profile/export-zip` | Export as portable offline ZIP bundle |
| `POST` | `/api/profile/diff` | Calculate diff between current state and a profile |
| `POST` | `/api/profile/export-batch` | Batch install from a profile |
