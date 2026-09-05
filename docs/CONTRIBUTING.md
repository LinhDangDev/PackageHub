# Contributing

## Setup

**Prerequisites:** Go 1.22+, Git, Windows 10/11 (for WebView2 testing).

```bash
git clone https://github.com/your-username/PackageHub.git
cd PackageHub
go mod tidy
```

## Build

```bash
go build -ldflags="-H windowsgui" -o packetinstall.exe ./cmd/packetinstall
```

## Test

```bash
go test -v ./...
```

All tests must pass before submitting a PR. Tests cover scanning, profile export/import, secret masking, and bundle generation.

## Pull Requests

1. Fork the repo and create your branch from `main`.
2. Make your changes and ensure `go test -v ./...` passes 100%.
3. Submit a PR with a clear description of changes.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Keep functions focused and testable.
- Add tests for new scanner or profile logic.
