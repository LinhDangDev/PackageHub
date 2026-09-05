package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

type FixRequest struct {
	Type        string `json:"type"`         // "tool", "path", "project"
	Manager     string `json:"manager"`      // "npm", "choco", "scoop"
	PackageName string `json:"package"`      // e.g. "@anthropic-ai/claude-code"
	ProjectPath string `json:"project_path"` // for project auto-fix
	IssueCode   string `json:"issue_code"`   // "OUTDATED", "PATH_MISSING", "AUDIT"
}

type FixResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Command string `json:"command"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

type FixDependencyRequest struct {
	ProjectPath string `json:"project_path"`
	Language    string `json:"language"`
	PackageName string `json:"package_name"`
	IsDev       bool   `json:"is_dev"`
}

// FixSingleDependency upgrades, installs, or pins a single specific dependency in a project.
func FixSingleDependency(req FixDependencyRequest) *FixResult {
	if req.ProjectPath == "" || req.PackageName == "" {
		return &FixResult{Success: false, Error: "project path and package name are required"}
	}

	var cmdStr string
	var cmd *exec.Cmd

	switch strings.ToLower(req.Language) {
	case "typescript", "javascript", "node.js":
		flag := "--save"
		if req.IsDev {
			flag = "--save-dev"
		}
		// Pin to latest resolved semver
		cmdStr = fmt.Sprintf("npm install %s@latest %s", req.PackageName, flag)
		cmd = exec.Command("cmd", "/c", cmdStr)
		cmd.Dir = req.ProjectPath
	case "go":
		cmdStr = fmt.Sprintf("go get -u %s@latest && go mod tidy", req.PackageName)
		cmd = exec.Command("cmd", "/c", cmdStr)
		cmd.Dir = req.ProjectPath
	case "python":
		cmdStr = fmt.Sprintf("pip install --upgrade %s", req.PackageName)
		cmd = exec.Command("cmd", "/c", cmdStr)
		cmd.Dir = req.ProjectPath
	case "rust":
		cmdStr = fmt.Sprintf("cargo update -p %s", req.PackageName)
		cmd = exec.Command("cmd", "/c", cmdStr)
		cmd.Dir = req.ProjectPath
	default:
		cmdStr = fmt.Sprintf("npm install %s@latest", req.PackageName)
		cmd = exec.Command("cmd", "/c", cmdStr)
		cmd.Dir = req.ProjectPath
	}

	hideConsoleFix(cmd)
	out, err := cmd.CombinedOutput()
	res := &FixResult{
		Success: err == nil,
		Command: cmdStr,
		Output:  string(out),
		Message: fmt.Sprintf("Successfully updated '%s' to latest resolved version in %s", req.PackageName, filepath.Base(req.ProjectPath)),
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

func hideConsoleFix(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
			HideWindow:    true,
		}
	}
}

// FixIssue automatically diagnoses and applies the exact remediation command for an issue.
func FixIssue(req FixRequest) *FixResult {
	switch req.Type {
	case "path":
		return fixPathIssue(req.Manager)
	case "project":
		return fixProjectIssues(req.ProjectPath)
	case "tool":
		return fixToolIssue(req.Manager, req.PackageName)
	default:
		return &FixResult{Success: false, Error: "unknown fix type"}
	}
}

func fixToolIssue(manager, pkgName string) *FixResult {
	var cmd *exec.Cmd
	var cmdStr string
	switch strings.ToLower(manager) {
	case "npm":
		cmdStr = fmt.Sprintf("npm install -g %s@latest", pkgName)
		cmd = exec.Command("cmd", "/c", cmdStr)
	case "choco":
		cmdStr = fmt.Sprintf("choco upgrade -y %s", pkgName)
		cmd = exec.Command("cmd", "/c", cmdStr)
	case "scoop":
		cmdStr = fmt.Sprintf("scoop update %s", pkgName)
		cmd = exec.Command("cmd", "/c", cmdStr)
	default:
		return &FixResult{Success: false, Error: "unsupported manager: " + manager}
	}

	hideConsoleFix(cmd)
	out, err := cmd.CombinedOutput()
	res := &FixResult{
		Success: err == nil,
		Command: cmdStr,
		Output:  string(out),
		Message: fmt.Sprintf("Applied upgrade for %s", pkgName),
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

func fixProjectIssues(projPath string) *FixResult {
	if projPath == "" {
		return &FixResult{Success: false, Error: "project path cannot be empty"}
	}

	// Check if package.json exists
	if _, err := os.Stat(filepath.Join(projPath, "package.json")); err == nil {
		cmdStr := "npm audit fix && npm update"
		cmd := exec.Command("cmd", "/c", "npm audit fix && npm update")
		cmd.Dir = projPath
		hideConsoleFix(cmd)
		out, err := cmd.CombinedOutput()
		res := &FixResult{
			Success: err == nil,
			Command: cmdStr,
			Output:  string(out),
			Message: "Executed npm audit fix and npm update",
		}
		if err != nil {
			res.Error = err.Error()
		}
		return res
	}

	// Check if go.mod exists
	if _, err := os.Stat(filepath.Join(projPath, "go.mod")); err == nil {
		cmdStr := "go get -u ./... && go mod tidy"
		cmd := exec.Command("cmd", "/c", "go get -u ./... && go mod tidy")
		cmd.Dir = projPath
		hideConsoleFix(cmd)
		out, err := cmd.CombinedOutput()
		res := &FixResult{
			Success: err == nil,
			Command: cmdStr,
			Output:  string(out),
			Message: "Updated Go modules and tidied go.mod",
		}
		if err != nil {
			res.Error = err.Error()
		}
		return res
	}

	return &FixResult{Success: false, Error: "no recognized project package manifest in: " + projPath}
}

func fixPathIssue(manager string) *FixResult {
	var missingDir string
	userProfile := os.Getenv("USERPROFILE")
	appData := os.Getenv("APPDATA")

	switch strings.ToLower(manager) {
	case "npm":
		missingDir = filepath.Join(appData, "npm")
	case "choco":
		missingDir = `C:\ProgramData\chocolatey\bin`
	case "scoop":
		missingDir = filepath.Join(userProfile, "scoop", "shims")
	}

	if missingDir == "" {
		return &FixResult{Success: false, Error: "no directory to add for manager: " + manager}
	}

	psScript := fmt.Sprintf(`
		$dir = '%s'
		$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
		if ($userPath -notlike "*$dir*") {
			[Environment]::SetEnvironmentVariable('Path', "$userPath;$dir", 'User')
			Write-Output "Successfully added $dir to User PATH"
		} else {
			Write-Output "$dir is already in User PATH"
		}
	`, missingDir)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	hideConsoleFix(cmd)
	out, err := cmd.CombinedOutput()

	res := &FixResult{
		Success: err == nil,
		Command: fmt.Sprintf("Add '%s' to User PATH", missingDir),
		Output:  string(out),
		Message: fmt.Sprintf("Added '%s' to your Windows User PATH environment variable!", missingDir),
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}
