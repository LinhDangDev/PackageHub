package profile

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
	"packetinstall/internal/model"
)

type SelectiveExportRequest struct {
	ProfileName   string   `json:"profile_name"`
	SelectedTools []string `json:"selected_tools"` // Package names
	SelectedSkills []string `json:"selected_skills"` // Skill names
	SelectedMcp   []string `json:"selected_mcp"`   // MCP server names
}

// FilterState creates a filtered SystemState containing only user-selected items.
func FilterState(state *model.SystemState, req SelectiveExportRequest) *model.SystemState {
	filtered := &model.SystemState{
		Packages:   make([]model.Package, 0),
		Skills:     make([]model.Skill, 0),
		McpServers: make([]model.McpServer, 0),
		ScannedAt:  state.ScannedAt,
		DurationMs: state.DurationMs,
	}

	toolSet := make(map[string]bool)
	for _, t := range req.SelectedTools {
		toolSet[strings.ToLower(t)] = true
	}

	skillSet := make(map[string]bool)
	for _, s := range req.SelectedSkills {
		skillSet[strings.ToLower(s)] = true
	}

	mcpSet := make(map[string]bool)
	for _, m := range req.SelectedMcp {
		mcpSet[strings.ToLower(m)] = true
	}

	// Filter Packages (if empty selection, include all)
	for _, p := range state.Packages {
		if len(toolSet) == 0 || toolSet[strings.ToLower(p.Name)] {
			filtered.Packages = append(filtered.Packages, p)
		}
	}

	// Filter Skills (if empty selection, include all)
	for _, s := range state.Skills {
		if len(skillSet) == 0 || skillSet[strings.ToLower(s.Name)] {
			filtered.Skills = append(filtered.Skills, s)
		}
	}

	// Filter MCP (if empty selection, include all)
	for _, m := range state.McpServers {
		if len(mcpSet) == 0 || mcpSet[strings.ToLower(m.Name)] {
			filtered.McpServers = append(filtered.McpServers, m)
		}
	}

	return filtered
}

// ExportZipBundle packages the YAML manifest, local skill files, and an auto-installer script into an offline ZIP bundle.
func ExportZipBundle(state *model.SystemState, req SelectiveExportRequest) ([]byte, error) {
	filtered := FilterState(state, req)
	profileName := req.ProfileName
	if profileName == "" {
		profileName = "portable-workstation"
	}

	yamlBytes, err := ExportProfileYAML(filtered, profileName)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// 1. Add packetinstall.yaml
	wYaml, err := zw.Create("packetinstall.yaml")
	if err != nil {
		return nil, err
	}
	if _, err := wYaml.Write(yamlBytes); err != nil {
		return nil, err
	}

	// 2. Generate install.ps1 auto-installer script
	p, _ := ExportProfile(filtered, profileName)
	emptyState := &model.SystemState{}
	diff := CalculateDiff(emptyState, p)
	cmds := GenerateInstallPlan(diff, runtime.GOOS, "./skills")

	var ps1Content strings.Builder
	ps1Content.WriteString("#!/usr/bin/env pwsh\n")
	ps1Content.WriteString("Write-Host '⚡ Installing packetinstall portable workstation bundle...' -ForegroundColor Cyan\n\n")
	for _, cmd := range cmds {
		ps1Content.WriteString(fmt.Sprintf("Write-Host 'Executing: %s' -ForegroundColor DarkGray\n", cmd))
		ps1Content.WriteString(cmd + "\n")
	}
	ps1Content.WriteString("\n# Copy packaged offline skills to target directory\n")
	ps1Content.WriteString("$targetSkillsDir = Join-Path $env:USERPROFILE '.agent' 'skills'\n")
	ps1Content.WriteString("if (!(Test-Path $targetSkillsDir)) { New-Item -ItemType Directory -Path $targetSkillsDir -Force }\n")
	ps1Content.WriteString("if (Test-Path './skills') { Copy-Item -Path './skills/*' -Destination $targetSkillsDir -Recurse -Force }\n")
	ps1Content.WriteString("Write-Host '✅ Workstation setup completed successfully!' -ForegroundColor Green\n")

	wPs1, err := zw.Create("install.ps1")
	if err == nil {
		_, _ = wPs1.Write([]byte(ps1Content.String()))
	}

	// 3. Add README.txt
	wReadme, err := zw.Create("README.txt")
	if err == nil {
		readmeText := fmt.Sprintf(`packetinstall Portable Workstation Bundle
===========================================
Profile: %s
Generated: %s

HOW TO USE ON A NEW MACHINE:
1. Extract this zip archive into any folder.
2. Open PowerShell as Administrator (if installing Chocolatey/Winget packages).
3. Run: .\install.ps1
4. All selected tools, CLIs, and AI Agent Skills will be automatically installed!
`, profileName, state.ScannedAt.Format("2006-01-02 15:04:05"))
		_, _ = wReadme.Write([]byte(readmeText))
	}

	// 4. Package local skill files into skills/ folder inside the zip
	for _, s := range filtered.Skills {
		if s.Path != "" {
			skillDirName := filepath.Base(s.Path)
			skillMdPath := filepath.Join(s.Path, "SKILL.md")
			if data, err := os.ReadFile(skillMdPath); err == nil {
				wSkill, err := zw.Create(fmt.Sprintf("skills/%s/SKILL.md", skillDirName))
				if err == nil {
					_, _ = wSkill.Write(data)
				}
			}
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportCustomYAML exports YAML with selected items.
func ExportCustomYAML(state *model.SystemState, req SelectiveExportRequest) ([]byte, error) {
	filtered := FilterState(state, req)
	profileName := req.ProfileName
	if profileName == "" {
		profileName = "custom-profile"
	}
	p, err := ExportProfile(filtered, profileName)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(p)
}
