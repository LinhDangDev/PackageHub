# ⚡ PackageHub (`packetinstall`)

> **Open-Source Developer Workstation, AI Agent Skills, Project Auditor & Offline ZIP Bundler**  
> An ultra-fast, zero-bloat developer environment manager written in Go with an embedded native Windows WebView2 desktop interface.

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20Native%20Desktop%20%7C%20CLI-lightgrey.svg)]()
[![GitHub Repo](https://img.shields.io/badge/github-LinhDangDev%2FPackageHub-cyan.svg)](https://github.com/LinhDangDev/PackageHub)

---

## 🌟 Highlights & Key Features

### 1. 🧠 AI Agent Skills Hub (Tách Biệt Hệ Sinh Thái & Tra Cứu Flags Chuẩn VividKit)
- **Phân chia nguồn gốc rõ ràng (Tool Origin)**: Lọc riêng theo **ClaudeKit (`ck:`)**, **OpenCode / OMP**, **Codex**, hoặc **Community / Universal**.
- **Tổ chức theo Cấp độ Kỹ năng (VividKit Skill Tiers)**:
  - 🟢 **Cơ bản (Beginner)**: `/ck:ask`, `/ck:brainstorm`, `/ck:fix`, `/ck:docs`, `/ck:interview-docs`...
  - 🟡 **Nâng cao (Intermediate)**: `/ck:plan`, `/ck:cook`, `/ck:test`, `/ck:review-pr`, `/ck:git`...
  - 🔴 **Chuyên nghiệp (Pro)**: `/ck:bootstrap`, `/ck:scout`, `/ck:ship`, `/ck:watzup`, `/ck:team`, `/ck:vibe`, `/ck:ai-artist`...
- **Tra Cứu Flags & Ý Nghĩa Chi Tiết**:
  - Tự động phân tích các cờ `--fast`, `--mode`, `--provider`, `--skip`, `--tdd`, `--parallel`... kèm giải thích chi tiết tác dụng.
  - Nút **`+ Thêm cờ`**: 1-click tự động nối cờ vào câu lệnh mẫu và copy vào clipboard.
- **Ngăn Kéo Tài Liệu Tương Tác Chuẩn VividKit (Docs Drawer)**:
  - Click bất kỳ skill nào để mở ngăn kéo tài liệu chi tiết.
  - Tự động đóng/thu gọn mượt mà khi click ra vùng trống bên ngoài hoặc bấm `ESC`.

### 2. 📁 Drive & Project Dependency Auditor
- **Quét Mã Nguồn Xuyên Suốt Mọi Ổ Đĩa** (`D:\IdeaSideProject`, `D:\`, `C:\Users\Dev`).
- **Tự Động Nhận Diện Ngôn Ngữ & Framework**: **TypeScript / JavaScript** (Next.js, React, Turbo, Express), **Go** (Gin, Fiber), **Rust**, **Python**.
- **Bắt Lỗi Dependencies Lỗi Thời & Độc Hại**: Bắt các thư viện đã khai tử (`request`, `moment`, `uuid@3`) hoặc dùng cờ wildcard rủi ro (`latest`, `*`).
- **Nút "⚡ Fix Gói Này" Độc Lập**: Sửa và khóa chuẩn phiên bản cho **đúng duy nhất dependency được chỉ định** mà không làm xáo trộn các thư viện khác trong dự án.
- Chấm điểm sức khỏe (A+, B, D) kèm nút **"⚡ Auto-Fix All Project Dependencies"**.

### 3. ⚡ Quản Lý & Tự Động Sửa Lỗi Công Cụ (Auto-Fix & Package Manager)
- Nút **Fix**: Tự động nâng cấp bản mới hoặc **tự động chèn thư mục còn thiếu vào Windows User PATH Registry** (`HKCU\Environment\Path`).
- Nút **Remove 🗑️**: Gỡ sạch công cụ và nút **Version ⇅** cho phép nâng cấp hoặc hạ cấp về bất kỳ version nào.
- Cài đặt công cụ mới qua **NPM Global**, **Chocolatey**, **Scoop**.

### 4. 📦 Tùy Chọn Export & Đóng Gói Offline ZIP Bundle
- Tích chọn riêng công cụ, skill, MCP server cần xuất.
- Xuất ra file YAML hoặc file **ZIP Bundle Offline (`packetinstall-bundle.zip`)**: Đóng gói cả mã nguồn các skill và file script `install.ps1` tự động bung và cài đặt trên máy mới không có internet.

### 5. 🖥️ Giao Diện Windows Native Độc Lập & Cố Định (Sticky Sidebar)
- Thanh tiêu đề tối màu chuẩn Windows 11 DWM (`#0b1019`).
- **Zero-Console Popup**: Tuyệt đối không còn hiện tượng terminal đen chớp tắt.
- Thanh Sidebar cố định 100%, chỉ cuộn phần nội dung bên phải.

---

## 🚀 Quick Start

### 1. Khởi chạy Ứng dụng Desktop Windows Native
Click đúp file `packetinstall.exe` trong File Explorer hoặc chạy trong PowerShell:
```powershell
.\packetinstall.exe
```

### 2. Dùng dòng lệnh CLI
```powershell
.\packetinstall.exe scan           # Quét nhanh tools & skills
.\packetinstall.exe audit          # Kiểm tra EOL & cập nhật
.\packetinstall.exe export -o dev.yaml # Xuất profile
.\packetinstall.exe apply dev.yaml --dry-run
```

---

## 📚 Documentation

- [Architecture & Engineering Design](docs/ARCHITECTURE.md)
- [User Guide & Operations Manual](docs/USAGE.md)
- [REST API Reference](docs/API.md)
- [Contributing Guide](docs/CONTRIBUTING.md)

---

## 🛠️ Architecture

```
PackageHub/
├── cmd/packetinstall/main.go     # CLI & Native Window entrypoint
├── internal/
│   ├── app/app.go                # Native Windows WebView2 Desktop Window & DWM Dark Mode
│   ├── model/types.go            # Domain models (Package, Project, Skill, McpServer, Profile)
│   ├── scanner/
│   │   ├── choco.go              # Chocolatey XML parser
│   │   ├── scoop.go              # Scoop apps parser
│   │   ├── npm.go                # Global NPM package.json parser
│   │   ├── skills.go             # Skills scanner with tool origin, tiers & flag parser
│   │   ├── mcp.go                # Claude Desktop & Cursor MCP parser
│   │   └── project.go            # Drive & Codebase dependency auditor
│   ├── installer/
│   │   ├── installer.go          # Install, Uninstall, Switch Version, Batch executor
│   │   └── fixer.go              # 1-Click Auto-Fix engine (Tools, PATH, Single Dependency)
│   ├── auditor/                  # EOL & Upstream registry update auditor
│   ├── profile/                  # YAML profile export/import & offline ZIP bundler
│   └── web/                      # HTTP REST API & Cyber-Industrial UI
├── docs/                         # Full documentation suite
└── plans/                        # Architectural research & TDD roadmap
```

---

## 📄 License

MIT License. Open-source và hoàn toàn miễn phí cho cộng đồng lập trình viên.
