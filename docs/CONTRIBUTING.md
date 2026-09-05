# 🤝 Contributing to PackageHub

Thank you for your interest in contributing to **PackageHub**!

---

## 1. Development Setup

### Prerequisites
- **Go 1.22+**
- **Git**
- **Windows 10/11** (for WebView2 testing)

### Clone & Build
```bash
git clone https://github.com/LinhDangDev/PackageHub.git
cd PackageHub

# Download dependencies
go mod download

# Build binary
go build -o packetinstall.exe ./cmd/packetinstall
```

---

## 2. Test-Driven Development (TDD) Workflow

Every feature or bug fix must follow strict TDD:
1. **Red**: Write unit tests in `*_test.go` asserting expected behaviors.
2. **Green**: Implement the minimal production code in `internal/*`.
3. **Refactor**: Clean up and optimize.

Run tests:
```bash
go test -v ./...
```

---

## 3. Architecture Rules
- **No Console Popups**: Never spawn unmanaged `exec.Command` on Windows without `CREATE_NO_WINDOW` (or `hideConsole`).
- **Direct Filesystem First**: Always prefer reading manifests (`package.json`, `.nuspec`, `.git/config`, `go.mod`) on disk rather than spawning slow CLI processes.
- **Secret Isolation**: Never serialize raw API credentials in exported profiles. Always mask with `${ENV_VAR}`.

---

## 4. Pull Request Process
1. Fork the repo and create your branch from `main`.
2. Ensure `go test -v ./...` passes 100%.
3. Submit your PR with a clear description of changes.
