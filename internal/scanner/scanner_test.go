package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanChoco(t *testing.T) {
	tmpDir := t.TempDir()
	pkgDir := filepath.Join(tmpDir, "ripgrep")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	nuspecContent := `<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2011/08/nuspec.xsd">
  <metadata>
    <id>ripgrep</id>
    <version>14.1.0</version>
    <title>ripgrep</title>
    <description>Fast search tool</description>
  </metadata>
</package>`
	if err := os.WriteFile(filepath.Join(pkgDir, "ripgrep.nuspec"), []byte(nuspecContent), 0644); err != nil {
		t.Fatal(err)
	}

	pkgs, err := ScanChoco(tmpDir)
	if err != nil {
		t.Fatalf("ScanChoco failed: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Name != "ripgrep" || pkgs[0].Version != "14.1.0" || pkgs[0].Manager != "choco" {
		t.Errorf("unexpected package: %+v", pkgs[0])
	}
}

func TestScanScoop(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "fzf")
	versionDir := filepath.Join(appDir, "0.54.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}

	pkgs, err := ScanScoop(tmpDir)
	if err != nil {
		t.Fatalf("ScanScoop failed: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Name != "fzf" || pkgs[0].Version != "0.54.0" || pkgs[0].Manager != "scoop" {
		t.Errorf("unexpected package: %+v", pkgs[0])
	}
}

func TestScanNpm(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Regular package
	pnpmDir := filepath.Join(tmpDir, "pnpm")
	if err := os.MkdirAll(pnpmDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pnpmDir, "package.json"), []byte(`{"name":"pnpm","version":"9.1.0"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Scoped package
	scopedDir := filepath.Join(tmpDir, "@anthropic-ai", "claude-code")
	if err := os.MkdirAll(scopedDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scopedDir, "package.json"), []byte(`{"name":"@anthropic-ai/claude-code","version":"2.1.261","description":"Claude CLI"}`), 0644); err != nil {
		t.Fatal(err)
	}

	pkgs, err := ScanNpm(tmpDir)
	if err != nil {
		t.Fatalf("ScanNpm failed: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}

	foundClaude := false
	for _, p := range pkgs {
		if p.Name == "@anthropic-ai/claude-code" {
			foundClaude = true
			if p.Version != "2.1.261" || p.Description != "Claude CLI" {
				t.Errorf("unexpected claude package details: %+v", p)
			}
		}
	}
	if !foundClaude {
		t.Errorf("expected to find @anthropic-ai/claude-code in scanned packages")
	}
}

func TestScanSkills(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "fable-thinking")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillContent := `---
name: fable-thinking
description: Disciplined reasoning protocol for code reviews and architecture.
---
# Fable Thinking
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatal(err)
	}

	skills, err := ScanSkills(tmpDir)
	if err != nil {
		t.Fatalf("ScanSkills failed: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "fable-thinking" {
		t.Errorf("expected skill name 'fable-thinking', got '%s'", skills[0].Name)
	}
	if skills[0].Description != "Disciplined reasoning protocol for code reviews and architecture." {
		t.Errorf("unexpected description: '%s'", skills[0].Description)
	}
}

func TestScanMcp(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "claude_desktop_config.json")
	mcpConfig := `{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "D:/Projects"],
      "env": {
        "DEBUG": "true"
      }
    }
  }
}`
	if err := os.WriteFile(tmpFile, []byte(mcpConfig), 0644); err != nil {
		t.Fatal(err)
	}

	servers, err := ScanMcpConfigFile(tmpFile, "claude-desktop")
	if err != nil {
		t.Fatalf("ScanMcpConfigFile failed: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 MCP server, got %d", len(servers))
	}
	if servers[0].Name != "filesystem" || servers[0].Command != "npx" || servers[0].Source != "claude-desktop" {
		t.Errorf("unexpected MCP server: %+v", servers[0])
	}
	if servers[0].Env["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true env var, got: %v", servers[0].Env)
	}
}

func TestScanAllCoordinator(t *testing.T) {
	tmpDir := t.TempDir()

	// Choco
	chocoDir := filepath.Join(tmpDir, "choco", "git")
	_ = os.MkdirAll(chocoDir, 0755)
	_ = os.WriteFile(filepath.Join(chocoDir, "git.nuspec"), []byte(`<package><metadata><id>git</id><version>2.43.0</version></metadata></package>`), 0644)

	// Scoop
	scoopDir := filepath.Join(tmpDir, "scoop", "neovim", "0.10.0")
	_ = os.MkdirAll(scoopDir, 0755)

	opts := ScanOptions{
		ChocoDir: filepath.Join(tmpDir, "choco"),
		ScoopDir: filepath.Join(tmpDir, "scoop"),
	}

	state, err := ScanAll(opts)
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}
	if len(state.Packages) != 2 {
		t.Errorf("expected 2 packages, got %d", len(state.Packages))
	}
	if state.DurationMs < 0 {
		t.Errorf("expected positive duration ms")
	}
}
