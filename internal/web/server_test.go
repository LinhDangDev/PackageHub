package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"packetinstall/internal/model"
	"packetinstall/internal/scanner"
)

func TestWebServerEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	chocoDir := filepath.Join(tmpDir, "choco", "git")
	_ = os.MkdirAll(chocoDir, 0755)
	_ = os.WriteFile(filepath.Join(chocoDir, "git.nuspec"), []byte(`<package><metadata><id>git</id><version>2.43.0</version></metadata></package>`), 0644)

	opts := scanner.ScanOptions{
		ChocoDir: filepath.Join(tmpDir, "choco"),
	}

	server := NewServer(opts)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	// 1. Test GET /api/scan
	resp, err := http.Get(ts.URL + "/api/scan")
	if err != nil {
		t.Fatalf("GET /api/scan failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got: %d", resp.StatusCode)
	}

	var state model.SystemState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if len(state.Packages) != 1 || state.Packages[0].Name != "git" {
		t.Errorf("unexpected scanned packages: %+v", state.Packages)
	}

	// 2. Test GET / (HTML Dashboard)
	respHtml, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer respHtml.Body.Close()

	if respHtml.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK for /, got: %d", respHtml.StatusCode)
	}

	// 3. Test POST /api/profile/diff
	yamlProfile := `
schema_version: "1.0"
metadata:
  name: test-profile
system_packages:
  - id: ripgrep
    windows:
      manager: choco
      package_id: ripgrep
`
	respDiff, err := http.Post(ts.URL+"/api/profile/diff", "application/x-yaml", strings.NewReader(yamlProfile))
	if err != nil {
		t.Fatalf("POST /api/profile/diff failed: %v", err)
	}
	defer respDiff.Body.Close()

	if respDiff.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK for diff, got: %d", respDiff.StatusCode)
	}

	// 4. Test POST /api/package/diagnose
	pkgJson := `{"manager":"npm","name":"@anthropic-ai/claude-code","version":"2.1.246","path":"."}`
	respDiag, err := http.Post(ts.URL+"/api/package/diagnose", "application/json", strings.NewReader(pkgJson))
	if err != nil {
		t.Fatalf("POST /api/package/diagnose failed: %v", err)
	}
	defer respDiag.Body.Close()
	if respDiag.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK for diagnose, got: %d", respDiag.StatusCode)
	}
	var diag model.PackageDiagnostics
	if err := json.NewDecoder(respDiag.Body).Decode(&diag); err != nil {
		t.Fatalf("failed to decode diagnostics response: %v", err)
	}
	if diag.Status == "" {
		t.Errorf("expected non-empty diagnostic status")
	}

	// 5. Test GET /api/cleaner/leftovers
	respLeftovers, err := http.Get(ts.URL + "/api/cleaner/leftovers?tool=ripgrep")
	if err != nil {
		t.Fatalf("GET /api/cleaner/leftovers failed: %v", err)
	}
	defer respLeftovers.Body.Close()
	if respLeftovers.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for leftovers, got %d", respLeftovers.StatusCode)
	}
	var leftoverRep model.LeftoverReport
	if err := json.NewDecoder(respLeftovers.Body).Decode(&leftoverRep); err != nil {
		t.Fatalf("failed to decode leftovers: %v", err)
	}
	if leftoverRep.ToolName != "ripgrep" {
		t.Errorf("expected tool_name 'ripgrep', got '%s'", leftoverRep.ToolName)
	}

	// 6. Test GET /api/cleaner/path
	respPath, err := http.Get(ts.URL + "/api/cleaner/path")
	if err != nil {
		t.Fatalf("GET /api/cleaner/path failed: %v", err)
	}
	defer respPath.Body.Close()
	if respPath.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for path audit, got %d", respPath.StatusCode)
	}

	// 7. Test GET /api/cleaner/caches
	respCaches, err := http.Get(ts.URL + "/api/cleaner/caches")
	if err != nil {
		t.Fatalf("GET /api/cleaner/caches failed: %v", err)
	}
	defer respCaches.Body.Close()
	if respCaches.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for caches, got %d", respCaches.StatusCode)
	}

	// 8. Test GET /api/cleaner/ports
	respPorts, err := http.Get(ts.URL + "/api/cleaner/ports")
	if err != nil {
		t.Fatalf("GET /api/cleaner/ports failed: %v", err)
	}
	defer respPorts.Body.Close()
	if respPorts.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for ports, got %d", respPorts.StatusCode)
	}
}
