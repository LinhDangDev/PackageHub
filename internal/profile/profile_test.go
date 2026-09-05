package profile

import (
	"strings"
	"testing"

	"packetinstall/internal/model"
)

func TestExportProfile_SecretMasking(t *testing.T) {
	state := &model.SystemState{
		Packages: []model.Package{
			{Manager: "choco", Name: "git", Version: "2.43.0"},
			{Manager: "npm", Name: "@anthropic-ai/claude-code", Version: "2.1.261"},
		},
		Skills: []model.Skill{
			{Name: "sequential-thinking", GitRemote: "https://github.com/example/sequential-thinking.git"},
		},
		McpServers: []model.McpServer{
			{
				Name:    "github",
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-github"},
				Env: map[string]string{
					"GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_super_secret_token_12345",
					"DEBUG":                        "true",
				},
			},
		},
	}

	yamlBytes, err := ExportProfileYAML(state, "dev-workstation")
	if err != nil {
		t.Fatalf("ExportProfileYAML failed: %v", err)
	}

	yamlStr := string(yamlBytes)

	// Invariant: Raw token must NEVER appear in output
	if strings.Contains(yamlStr, "ghp_super_secret_token_12345") {
		t.Errorf("CRITICAL SECURITY VIOLATION: raw secret token found in exported YAML!")
	}

	// Must contain masked placeholder
	if !strings.Contains(yamlStr, "${GITHUB_PERSONAL_ACCESS_TOKEN}") {
		t.Errorf("expected masked placeholder '${GITHUB_PERSONAL_ACCESS_TOKEN}', got:\n%s", yamlStr)
	}

	// Non-secret env var DEBUG should remain intact
	if !strings.Contains(yamlStr, "DEBUG: \"true\"") && !strings.Contains(yamlStr, "DEBUG: true") {
		t.Errorf("expected non-secret env var DEBUG to remain, got:\n%s", yamlStr)
	}
}

func TestCalculateDiff(t *testing.T) {
	currentState := &model.SystemState{
		Packages: []model.Package{
			{Manager: "choco", Name: "git", Version: "2.43.0"},
		},
		Skills: []model.Skill{
			{Name: "sequential-thinking"},
		},
		McpServers: []model.McpServer{},
	}

	profile := &model.Profile{
		SchemaVersion: "1.0",
		GlobalCLIs: map[string][]string{
			"npm": {"@anthropic-ai/claude-code", "pnpm"},
		},
		SystemTools: []model.SystemPackageSpec{
			{
				ID: "git",
				Windows: &model.PlatformPackageSpec{
					Manager:   "choco",
					PackageID: "git",
				},
			},
			{
				ID: "ripgrep",
				Windows: &model.PlatformPackageSpec{
					Manager:   "choco",
					PackageID: "ripgrep",
				},
			},
		},
		Skills: []model.SkillProfileSpec{
			{Name: "sequential-thinking"},
			{Name: "aesthetic", RepoURL: "https://github.com/example/aesthetic.git"},
		},
		McpServers: []model.McpServerProfileSpec{
			{Name: "filesystem", Command: "npx"},
		},
	}

	diff := CalculateDiff(currentState, profile)

	if len(diff.AlreadyInstalled) != 1 || diff.AlreadyInstalled[0] != "git" {
		t.Errorf("expected 'git' in AlreadyInstalled, got: %v", diff.AlreadyInstalled)
	}

	if len(diff.PendingSystemPackages) != 1 || diff.PendingSystemPackages[0].ID != "ripgrep" {
		t.Errorf("expected 'ripgrep' in PendingSystemPackages, got: %v", diff.PendingSystemPackages)
	}

	if len(diff.PendingGlobalCLIs["npm"]) != 2 {
		t.Errorf("expected 2 npm CLIs pending install, got: %v", diff.PendingGlobalCLIs["npm"])
	}

	if len(diff.MissingSkills) != 1 || diff.MissingSkills[0].Name != "aesthetic" {
		t.Errorf("expected 'aesthetic' in MissingSkills, got: %v", diff.MissingSkills)
	}

	if len(diff.MissingMcpServers) != 1 || diff.MissingMcpServers[0].Name != "filesystem" {
		t.Errorf("expected 'filesystem' in MissingMcpServers, got: %v", diff.MissingMcpServers)
	}
}

func TestGenerateInstallPlan(t *testing.T) {
	diff := &ProfileDiff{
		PendingSystemPackages: []model.SystemPackageSpec{
			{
				ID: "ripgrep",
				Windows: &model.PlatformPackageSpec{
					Manager:   "choco",
					PackageID: "ripgrep",
				},
			},
		},
		PendingGlobalCLIs: map[string][]string{
			"npm": {"@anthropic-ai/claude-code"},
		},
		MissingSkills: []model.SkillProfileSpec{
			{
				Name:    "aesthetic",
				RepoURL: "https://github.com/example/aesthetic.git",
			},
		},
	}

	commands := GenerateInstallPlan(diff, "windows", "C:/Users/Dev/.agent/skills")
	if len(commands) < 3 {
		t.Fatalf("expected at least 3 install commands, got %d: %v", len(commands), commands)
	}

	hasChoco := false
	hasNpm := false
	hasGit := false

	for _, cmd := range commands {
		if strings.Contains(cmd, "choco install -y ripgrep") {
			hasChoco = true
		}
		if strings.Contains(cmd, "npm install -g @anthropic-ai/claude-code") {
			hasNpm = true
		}
		if strings.Contains(cmd, "git clone") && strings.Contains(cmd, "aesthetic") {
			hasGit = true
		}
	}

	if !hasChoco {
		t.Errorf("expected choco install command in plan, got: %v", commands)
	}
	if !hasNpm {
		t.Errorf("expected npm install command in plan, got: %v", commands)
	}
	if !hasGit {
		t.Errorf("expected git clone command in plan, got: %v", commands)
	}
}
