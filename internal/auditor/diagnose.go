package auditor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"packetinstall/internal/model"
)

type detailedNpmPackageJson struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Bin         json.RawMessage   `json:"bin"`
	Engines     map[string]string `json:"engines"`
}

// DiagnosePackage inspects an installed package for updates, PATH availability, and engine issues.
func (a *AuditorClient) DiagnosePackage(pkg model.Package) *model.PackageDiagnostics {
	diag := &model.PackageDiagnostics{
		Status:    "HEALTHY",
		InPath:    true,
		Issues:    make([]model.PackageIssue, 0),
		CheckedAt: time.Now(),
	}

	switch pkg.Manager {
	case "npm":
		a.diagnoseNpmPackage(&pkg, diag)
	case "choco":
		a.diagnoseChocoPackage(&pkg, diag)
	default:
		diag.Status = "HEALTHY"
	}

	// Compute overall status based on issues
	for _, iss := range diag.Issues {
		if iss.Code == "EOL" {
			diag.Status = "EOL_CRITICAL"
			break
		}
		if iss.Severity == "ERROR" {
			diag.Status = "BROKEN"
			break
		}
		if iss.Code == "OUTDATED" && diag.Status != "BROKEN" && diag.Status != "EOL_CRITICAL" {
			diag.Status = "OUTDATED"
		} else if iss.Severity == "WARNING" && diag.Status == "HEALTHY" {
			diag.Status = "WARNING"
		}
	}

	return diag
}

func (a *AuditorClient) diagnoseNpmPackage(pkg *model.Package, diag *model.PackageDiagnostics) {
	pkgJsonPath := filepath.Join(pkg.Path, "package.json")
	data, err := os.ReadFile(pkgJsonPath)
	var pj detailedNpmPackageJson
	if err == nil {
		_ = json.Unmarshal(data, &pj)
	}

	// 1. Binary & PATH check
	binName := ""
	if len(pj.Bin) > 0 {
		var singleBin string
		var binMap map[string]string
		if err := json.Unmarshal(pj.Bin, &singleBin); err == nil {
			binName = pkg.Name
		} else if err := json.Unmarshal(pj.Bin, &binMap); err == nil {
			for k := range binMap {
				binName = k
				break
			}
		}
	} else {
		// Fallback for tools where binary matches package name suffix
		parts := strings.Split(pkg.Name, "/")
		binName = parts[len(parts)-1]
	}

	if binName != "" {
		diag.BinaryName = binName
		_, lookErr := exec.LookPath(binName)
		if lookErr != nil {
			// Check .cmd extension on Windows
			_, lookCmdErr := exec.LookPath(binName + ".cmd")
			if lookCmdErr != nil {
				diag.InPath = false
				diag.Issues = append(diag.Issues, model.PackageIssue{
					Severity: "WARNING",
					Code:     "PATH_MISSING",
					Message:  fmt.Sprintf("Binary '%s' is installed in node_modules but not found in current PATH.", binName),
					FixHint:  "Ensure your global npm prefix directory (e.g. %APPDATA%\\npm) is included in your system PATH environment variable.",
				})
			}
		}
	}

	// 2. Engines check
	if pj.Engines != nil {
		if nodeReq, ok := pj.Engines["node"]; ok {
			diag.EngineReq = nodeReq
		}
	}

	// 3. Upstream Registry Update Check
	if item, err := a.CheckNpmPackage(pkg.Name, pkg.Version); err == nil && item != nil {
		diag.LatestVersion = item.LatestVersion
		diag.UpdateCommand = item.UpdateCommand
		if item.Status == "UPDATE_AVAILABLE" {
			diag.Issues = append(diag.Issues, model.PackageIssue{
				Severity: "INFO",
				Code:     "OUTDATED",
				Message:  fmt.Sprintf("Update available: %s -> %s", pkg.Version, item.LatestVersion),
				FixHint:  item.UpdateCommand,
			})
		}
	}
}

func (a *AuditorClient) diagnoseChocoPackage(pkg *model.Package, diag *model.PackageDiagnostics) {
	diag.BinaryName = pkg.Name
	diag.UpdateCommand = fmt.Sprintf("choco upgrade -y %s", pkg.Name)

	// Check if executable exists in PATH for common command line tools
	commonTools := []string{"git", "ripgrep", "rg", "neovim", "nvim", "lazygit", "fd", "python"}
	isCommonTool := false
	for _, t := range commonTools {
		if strings.EqualFold(pkg.Name, t) {
			isCommonTool = true
			break
		}
	}

	if isCommonTool {
		_, lookErr := exec.LookPath(pkg.Name)
		if lookErr != nil {
			diag.InPath = false
			diag.Issues = append(diag.Issues, model.PackageIssue{
				Severity: "WARNING",
				Code:     "PATH_MISSING",
				Message:  fmt.Sprintf("Chocolatey tool '%s' was not found in active PATH.", pkg.Name),
				FixHint:  "Ensure 'C:\\ProgramData\\chocolatey\\bin' is in your system PATH, or restart your terminal.",
			})
		}
	}
}
