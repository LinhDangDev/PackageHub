package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanProjects_NodeAndGo(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Setup mock Next.js / TypeScript project
	nextDir := filepath.Join(tmpDir, "my-web-app")
	_ = os.MkdirAll(nextDir, 0755)
	_ = os.WriteFile(filepath.Join(nextDir, "tsconfig.json"), []byte("{}"), 0644)
	pkgJson := `{
  "name": "my-web-app",
  "dependencies": {
    "next": "15.1.0",
    "react": "19.0.0",
    "moment": "2.29.4"
  },
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}`
	_ = os.WriteFile(filepath.Join(nextDir, "package.json"), []byte(pkgJson), 0644)

	// 2. Setup mock Go project
	goDir := filepath.Join(tmpDir, "my-go-api")
	_ = os.MkdirAll(goDir, 0755)
	goMod := `module my-go-api

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
)
`
	_ = os.WriteFile(filepath.Join(goDir, "go.mod"), []byte(goMod), 0644)

	// Run scanner
	res, err := ScanProjects(tmpDir, 2)
	if err != nil {
		t.Fatalf("ScanProjects failed: %v", err)
	}

	if len(res.Projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(res.Projects))
	}

	var foundNext, foundGo bool
	for _, p := range res.Projects {
		if p.Name == "my-web-app" {
			foundNext = true
			if p.Language != "TypeScript" {
				t.Errorf("expected TypeScript, got %s", p.Language)
			}
			if p.Framework != "Next.js" {
				t.Errorf("expected Next.js, got %s", p.Framework)
			}
			if p.IssueCount != 1 {
				t.Errorf("expected 1 issue (moment deprecated), got %d", p.IssueCount)
			}
		}
		if p.Name == "my-go-api" {
			foundGo = true
			if p.Language != "Go" {
				t.Errorf("expected Go, got %s", p.Language)
			}
			if p.Framework != "Gin" {
				t.Errorf("expected Gin, got %s", p.Framework)
			}
		}
	}

	if !foundNext {
		t.Errorf("Next.js project not found in scan results")
	}
	if !foundGo {
		t.Errorf("Go project not found in scan results")
	}
}
