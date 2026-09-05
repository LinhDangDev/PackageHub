package scanner

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"packetinstall/internal/model"
)

// ScanOptions configures directories to scan. Empty values default to system locations.
type ScanOptions struct {
	ChocoDir     string
	ScoopDir     string
	NpmDir       string
	SkillsDirs   []string
	McpFiles     []McpScanTarget
}

type McpScanTarget struct {
	Path   string
	Source string
}

// DefaultScanOptions returns platform-specific default paths for scanning.
func DefaultScanOptions() ScanOptions {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		userProfile = os.Getenv("HOME")
	}
	appData := os.Getenv("APPDATA")

	opts := ScanOptions{
		ChocoDir: `C:\ProgramData\chocolatey\lib`,
		ScoopDir: filepath.Join(userProfile, "scoop", "apps"),
		NpmDir:   filepath.Join(appData, "npm", "node_modules"),
		SkillsDirs: []string{
			filepath.Join(userProfile, ".agent", "skills"),
			filepath.Join(userProfile, ".omp", "agent", "skills"),
			filepath.Join(userProfile, ".claude", "skills"),
		},
		McpFiles: []McpScanTarget{
			{
				Path:   filepath.Join(appData, "Claude", "claude_desktop_config.json"),
				Source: "claude-desktop",
			},
			{
				Path:   filepath.Join(appData, "Cursor", "User", "globalStorage", "mcp.json"),
				Source: "cursor",
			},
		},
	}

	return opts
}

// ScanAll orchestrates concurrent scanning across all package managers, skills, and MCP configs.
func ScanAll(opts ScanOptions) (*model.SystemState, error) {
	start := time.Now()
	state := &model.SystemState{
		Packages:   make([]model.Package, 0),
		Skills:     make([]model.Skill, 0),
		McpServers: make([]model.McpServer, 0),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	// 1. Choco
	if opts.ChocoDir != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pkgs, _ := ScanChoco(opts.ChocoDir)
			if len(pkgs) > 0 {
				mu.Lock()
				state.Packages = append(state.Packages, pkgs...)
				mu.Unlock()
			}
		}()
	}

	// 2. Scoop
	if opts.ScoopDir != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pkgs, _ := ScanScoop(opts.ScoopDir)
			if len(pkgs) > 0 {
				mu.Lock()
				state.Packages = append(state.Packages, pkgs...)
				mu.Unlock()
			}
		}()
	}

	// 3. NPM Global
	if opts.NpmDir != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pkgs, _ := ScanNpm(opts.NpmDir)
			if len(pkgs) > 0 {
				mu.Lock()
				state.Packages = append(state.Packages, pkgs...)
				mu.Unlock()
			}
		}()
	}

	// 4. Skills (multiple dirs)
	for _, sDir := range opts.SkillsDirs {
		dir := sDir
		wg.Add(1)
		go func() {
			defer wg.Done()
			skills, _ := ScanSkills(dir)
			if len(skills) > 0 {
				mu.Lock()
				state.Skills = append(state.Skills, skills...)
				mu.Unlock()
			}
		}()
	}

	// 5. MCP Files
	for _, mTarget := range opts.McpFiles {
		target := mTarget
		wg.Add(1)
		go func() {
			defer wg.Done()
			servers, _ := ScanMcpConfigFile(target.Path, target.Source)
			if len(servers) > 0 {
				mu.Lock()
				state.McpServers = append(state.McpServers, servers...)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	state.ScannedAt = time.Now()
	state.DurationMs = time.Since(start).Milliseconds()

	return state, nil
}
