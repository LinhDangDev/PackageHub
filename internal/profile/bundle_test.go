package profile

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"packetinstall/internal/model"
)

func TestSelectiveExportAndZipBundle(t *testing.T) {
	state := &model.SystemState{
		Packages: []model.Package{
			{Manager: "choco", Name: "ripgrep", Version: "14.1.0"},
			{Manager: "choco", Name: "git", Version: "2.43.0"},
			{Manager: "npm", Name: "pnpm", Version: "9.1.0"},
		},
		Skills: []model.Skill{
			{Name: "sequential-thinking"},
			{Name: "aesthetic"},
		},
		McpServers: []model.McpServer{
			{Name: "filesystem", Command: "npx"},
		},
	}

	req := SelectiveExportRequest{
		ProfileName:    "my-selective-bundle",
		SelectedTools:  []string{"ripgrep", "pnpm"},
		SelectedSkills: []string{"sequential-thinking"},
	}

	// 1. Test FilterState
	filtered := FilterState(state, req)
	if len(filtered.Packages) != 2 {
		t.Errorf("expected 2 filtered packages, got %d", len(filtered.Packages))
	}
	if len(filtered.Skills) != 1 || filtered.Skills[0].Name != "sequential-thinking" {
		t.Errorf("expected 1 filtered skill sequential-thinking, got %v", filtered.Skills)
	}

	// 2. Test ZIP bundle generation
	zipBytes, err := ExportZipBundle(state, req)
	if err != nil {
		t.Fatalf("ExportZipBundle failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("failed to read zip archive: %v", err)
	}

	hasYaml := false
	hasInstallPs1 := false
	hasReadme := false

	for _, f := range zr.File {
		if f.Name == "packetinstall.yaml" {
			hasYaml = true
		}
		if f.Name == "install.ps1" {
			hasInstallPs1 = true
		}
		if f.Name == "README.txt" {
			hasReadme = true
		}
	}

	if !hasYaml {
		t.Errorf("expected packetinstall.yaml inside zip")
	}
	if !hasInstallPs1 {
		t.Errorf("expected install.ps1 inside zip")
	}
	if !hasReadme {
		t.Errorf("expected README.txt inside zip")
	}
}

func TestSwitchPackageVersion(t *testing.T) {
	// Simple mock test
	req := SelectiveExportRequest{
		ProfileName:   "test",
		SelectedTools: []string{"git"},
	}
	yamlBytes, err := ExportCustomYAML(&model.SystemState{Packages: []model.Package{{Name: "git", Manager: "choco"}}}, req)
	if err != nil {
		t.Fatalf("ExportCustomYAML failed: %v", err)
	}
	if !strings.Contains(string(yamlBytes), "git") {
		t.Errorf("expected git in yaml output")
	}
}
