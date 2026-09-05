package profile

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"packetinstall/internal/model"
)

type ProfileDiff struct {
	AlreadyInstalled      []string
	PendingSystemPackages []model.SystemPackageSpec
	PendingGlobalCLIs     map[string][]string
	MissingSkills         []model.SkillProfileSpec
	MissingMcpServers     []model.McpServerProfileSpec
}

// ExportProfile builds a declarative Profile struct from SystemState with secrets masked.
func ExportProfile(state *model.SystemState, profileName string) (*model.Profile, error) {
	osName := runtime.GOOS
	if profileName == "" {
		profileName = "default-profile"
	}

	p := &model.Profile{
		SchemaVersion: "1.0",
		Metadata: model.ProfileMetadata{
			Name:      profileName,
			CreatedAt: time.Now().Format(time.RFC3339),
			SourceOS:  osName,
		},
		GlobalCLIs: make(map[string][]string),
	}

	// 1. Classify packages
	for _, pkg := range state.Packages {
		switch pkg.Manager {
		case "npm", "pipx", "cargo":
			p.GlobalCLIs[pkg.Manager] = append(p.GlobalCLIs[pkg.Manager], pkg.Name)
		case "choco", "scoop", "winget":
			spec := model.SystemPackageSpec{
				ID: pkg.Name,
			}
			if osName == "windows" {
				spec.Windows = &model.PlatformPackageSpec{
					Manager:   pkg.Manager,
					PackageID: pkg.Name,
				}
			}
			p.SystemTools = append(p.SystemTools, spec)
		}
	}

	// 2. Skills
	for _, s := range state.Skills {
		p.Skills = append(p.Skills, model.SkillProfileSpec{
			Name:      s.Name,
			RepoURL:   s.GitRemote,
			CommitSHA: s.CommitSHA,
		})
	}

	// 3. MCP Servers with Secret Masking
	for _, m := range state.McpServers {
		maskedEnv := make(map[string]string)
		for k, v := range m.Env {
			if isSecretKeyOrValue(k, v) {
				maskedEnv[k] = fmt.Sprintf("${%s}", k)
			} else {
				maskedEnv[k] = v
			}
		}

		p.McpServers = append(p.McpServers, model.McpServerProfileSpec{
			Name:    m.Name,
			Command: m.Command,
			Args:    m.Args,
			Env:     maskedEnv,
		})
	}

	return p, nil
}

// ExportProfileYAML exports state as YAML bytes.
func ExportProfileYAML(state *model.SystemState, profileName string) ([]byte, error) {
	p, err := ExportProfile(state, profileName)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(p)
}

// ImportProfileYAML parses a YAML manifest into a Profile struct.
func ImportProfileYAML(data []byte) (*model.Profile, error) {
	var p model.Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// CalculateDiff compares the target machine's current state with the desired profile.
func CalculateDiff(currentState *model.SystemState, targetProfile *model.Profile) *ProfileDiff {
	diff := &ProfileDiff{
		PendingGlobalCLIs: make(map[string][]string),
	}

	installedPkgMap := make(map[string]bool)
	for _, p := range currentState.Packages {
		installedPkgMap[strings.ToLower(p.Name)] = true
	}

	installedSkillMap := make(map[string]bool)
	for _, s := range currentState.Skills {
		installedSkillMap[strings.ToLower(s.Name)] = true
	}

	installedMcpMap := make(map[string]bool)
	for _, m := range currentState.McpServers {
		installedMcpMap[strings.ToLower(m.Name)] = true
	}

	// Diff System Packages
	for _, sp := range targetProfile.SystemTools {
		if installedPkgMap[strings.ToLower(sp.ID)] {
			diff.AlreadyInstalled = append(diff.AlreadyInstalled, sp.ID)
		} else {
			diff.PendingSystemPackages = append(diff.PendingSystemPackages, sp)
		}
	}

	// Diff Global CLIs
	for mgr, pkgs := range targetProfile.GlobalCLIs {
		for _, pkgName := range pkgs {
			if installedPkgMap[strings.ToLower(pkgName)] {
				diff.AlreadyInstalled = append(diff.AlreadyInstalled, pkgName)
			} else {
				diff.PendingGlobalCLIs[mgr] = append(diff.PendingGlobalCLIs[mgr], pkgName)
			}
		}
	}

	// Diff Skills
	for _, sk := range targetProfile.Skills {
		if !installedSkillMap[strings.ToLower(sk.Name)] {
			diff.MissingSkills = append(diff.MissingSkills, sk)
		}
	}

	// Diff MCP Servers
	for _, ms := range targetProfile.McpServers {
		if !installedMcpMap[strings.ToLower(ms.Name)] {
			diff.MissingMcpServers = append(diff.MissingMcpServers, ms)
		}
	}

	return diff
}

// GenerateInstallPlan turns a ProfileDiff into non-interactive execution commands.
func GenerateInstallPlan(diff *ProfileDiff, targetOS string, skillsTargetDir string) []string {
	var cmds []string

	// 1. System Packages
	for _, sp := range diff.PendingSystemPackages {
		if targetOS == "windows" && sp.Windows != nil {
			switch sp.Windows.Manager {
			case "choco":
				cmds = append(cmds, fmt.Sprintf("choco install -y %s", sp.Windows.PackageID))
			case "scoop":
				cmds = append(cmds, fmt.Sprintf("scoop install %s", sp.Windows.PackageID))
			case "winget":
				cmds = append(cmds, fmt.Sprintf("winget install --silent --accept-source-agreements --accept-package-agreements --id %s", sp.Windows.PackageID))
			}
		}
	}

	// 2. Global CLIs
	for mgr, pkgs := range diff.PendingGlobalCLIs {
		switch mgr {
		case "npm":
			for _, p := range pkgs {
				cmds = append(cmds, fmt.Sprintf("npm install -g %s", p))
			}
		case "pipx":
			for _, p := range pkgs {
				cmds = append(cmds, fmt.Sprintf("pipx install %s", p))
			}
		case "cargo":
			for _, p := range pkgs {
				cmds = append(cmds, fmt.Sprintf("cargo install %s", p))
			}
		}
	}

	// 3. Skills
	for _, sk := range diff.MissingSkills {
		if sk.RepoURL != "" {
			dest := filepath.Join(skillsTargetDir, sk.Name)
			cmds = append(cmds, fmt.Sprintf("git clone %s %s", sk.RepoURL, dest))
		}
	}

	return cmds
}

func isSecretKeyOrValue(key, val string) bool {
	upperKey := strings.ToUpper(key)
	secretIndicators := []string{"TOKEN", "KEY", "SECRET", "PASSWORD", "AUTH", "PASSWD", "CREDENTIAL"}
	for _, ind := range secretIndicators {
		if strings.Contains(upperKey, ind) {
			return true
		}
	}

	lowerVal := strings.ToLower(val)
	if strings.HasPrefix(lowerVal, "sk-") || strings.HasPrefix(lowerVal, "ghp_") || strings.HasPrefix(lowerVal, "bearer ") {
		return true
	}

	return false
}
