package cleaner

import (
	"os"
	"path/filepath"
	"testing"

	"packetinstall/internal/model"
)

func TestNormalizeToolName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@anthropic-ai/claude-code", "claude-code"},
		{"ripgrep", "ripgrep"},
		{"Docker Desktop", "docker desktop"},
		{"path/to/my-tool", "my-tool"},
		{"", ""},
	}

	for _, tt := range tests {
		got := normalizeToolName(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeToolName(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestScanAndPurgeLeftovers(t *testing.T) {
	tmpDir := t.TempDir()
	toolName := "dummy-test-tool"

	// Mock APPDATA to tmpDir
	origAppData := os.Getenv("APPDATA")
	defer os.Setenv("APPDATA", origAppData)
	os.Setenv("APPDATA", tmpDir)

	targetDir := filepath.Join(tmpDir, toolName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(targetDir, "config.json")
	if err := os.WriteFile(testFile, []byte(`{"installed": true}`), 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Scan Leftovers
	report := ScanLeftovers(toolName)
	if report.TotalItems == 0 {
		t.Fatalf("expected at least 1 leftover item for %s, got 0", toolName)
	}

	found := false
	for _, it := range report.Items {
		if it.Path == targetDir {
			found = true
			if it.Type != "dir" {
				t.Errorf("expected type 'dir', got %s", it.Type)
			}
			if it.Size == 0 {
				t.Errorf("expected size > 0, got %d", it.Size)
			}
		}
	}
	if !found {
		t.Errorf("targetDir %s not found in leftover report", targetDir)
	}

	// 2. Purge Leftovers
	purgeRes := PurgeLeftovers(model.PurgeLeftoversRequest{
		ToolName: toolName,
		ItemIDs:  []string{report.Items[0].ID},
	})
	if !purgeRes.Success || purgeRes.PurgedCount == 0 {
		t.Errorf("PurgeLeftovers failed: %+v", purgeRes)
	}

	// Verify directory is deleted
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("expected targetDir to be deleted, but still exists")
	}
}

func TestAuditPath(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentDir := filepath.Join(tmpDir, "does-not-exist-xyz-999")

	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+nonExistentDir)

	report := AuditPath()
	if report.TotalEntries < 2 {
		t.Fatalf("expected at least 2 path entries, got %d", report.TotalEntries)
	}

	foundValid := false
	foundDead := false

	for _, e := range report.Entries {
		if e.Path == tmpDir && e.Exists {
			foundValid = true
		}
		if e.Path == nonExistentDir && !e.Exists {
			foundDead = true
		}
	}

	if !foundValid {
		t.Errorf("expected tmpDir to be reported as valid")
	}
	if !foundDead {
		t.Errorf("expected nonExistentDir to be reported as dead/non-existent")
	}
	if report.DeadCount < 1 {
		t.Errorf("expected DeadCount >= 1, got %d", report.DeadCount)
	}
}

func TestAuditDevCaches(t *testing.T) {
	report := AuditDevCaches()
	if len(report.Caches) == 0 {
		t.Errorf("expected dev cache specs to be enumerated")
	}

	// Check if known IDs exist in list
	ids := make(map[string]bool)
	for _, c := range report.Caches {
		ids[c.ID] = true
	}

	if !ids["npm"] || !ids["go_build"] || !ids["pip"] {
		t.Errorf("expected npm, go_build, and pip in dev cache list")
	}
}

func TestAuditDevPorts(t *testing.T) {
	report := AuditDevPorts()
	// AuditDevPorts should succeed and return a valid report
	t.Logf("Total active listening dev ports found: %d", report.TotalListening)
	for _, p := range report.Ports {
		if p.Port <= 0 {
			t.Errorf("invalid port number: %d", p.Port)
		}
	}
}
