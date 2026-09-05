package installer

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"packetinstall/internal/model"
)

type InstallResult struct {
	Command string `json:"command"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

type UninstallRequest struct {
	Manager string `json:"manager"`
	Package string `json:"package"`
}

type SwitchVersionRequest struct {
	Manager string `json:"manager"`
	Package string `json:"package"`
	Version string `json:"version"`
	Global  bool   `json:"global"`
}

type BatchExecutionStep struct {
	Command string `json:"command"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

type BatchExecutionResult struct {
	TotalSteps     int                  `json:"total_steps"`
	CompletedSteps int                  `json:"completed_steps"`
	Success        bool                 `json:"success"`
	Steps          []BatchExecutionStep `json:"steps"`
}

func hideConsole(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
			HideWindow:    true,
		}
	}
}

// InstallPackage executes an automated installation command for a tool.
func InstallPackage(req model.InstallPackageRequest) *InstallResult {
	pkgName := strings.TrimSpace(req.Package)
	if pkgName == "" {
		return &InstallResult{
			Success: false,
			Error:   "package name cannot be empty",
		}
	}

	var cmdName string
	var args []string

	switch strings.ToLower(req.Manager) {
	case "npm":
		cmdName = "npm"
		args = []string{"install"}
		if req.Global {
			args = append(args, "-g")
		}
		args = append(args, pkgName)
	case "choco":
		cmdName = "choco"
		args = []string{"install", "-y", pkgName}
	case "scoop":
		cmdName = "scoop"
		args = []string{"install", pkgName}
	default:
		return &InstallResult{
			Success: false,
			Error:   fmt.Sprintf("unsupported package manager: %s", req.Manager),
		}
	}

	fullCmd := fmt.Sprintf("%s %s", cmdName, strings.Join(args, " "))
	cmd := exec.Command(cmdName, args...)
	hideConsole(cmd)
	out, err := cmd.CombinedOutput()

	res := &InstallResult{
		Command: fullCmd,
		Output:  string(out),
		Success: err == nil,
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

// UninstallPackage uninstalls a tool from NPM, Choco, or Scoop.
func UninstallPackage(req UninstallRequest) *InstallResult {
	pkgName := strings.TrimSpace(req.Package)
	if pkgName == "" {
		return &InstallResult{
			Success: false,
			Error:   "package name cannot be empty",
		}
	}

	var cmdName string
	var args []string

	switch strings.ToLower(req.Manager) {
	case "npm":
		cmdName = "npm"
		args = []string{"uninstall", "-g", pkgName}
	case "choco":
		cmdName = "choco"
		args = []string{"uninstall", "-y", pkgName}
	case "scoop":
		cmdName = "scoop"
		args = []string{"uninstall", pkgName}
	default:
		return &InstallResult{
			Success: false,
			Error:   fmt.Sprintf("unsupported package manager: %s", req.Manager),
		}
	}

	fullCmd := fmt.Sprintf("%s %s", cmdName, strings.Join(args, " "))
	cmd := exec.Command(cmdName, args...)
	hideConsole(cmd)
	out, err := cmd.CombinedOutput()

	res := &InstallResult{
		Command: fullCmd,
		Output:  string(out),
		Success: err == nil,
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

// SwitchPackageVersion upgrades or downgrades a package to a specific version.
func SwitchPackageVersion(req SwitchVersionRequest) *InstallResult {
	pkgName := strings.TrimSpace(req.Package)
	ver := strings.TrimSpace(req.Version)
	if pkgName == "" {
		return &InstallResult{Success: false, Error: "package name cannot be empty"}
	}
	if ver == "" {
		ver = "latest"
	}

	var cmdName string
	var args []string

	switch strings.ToLower(req.Manager) {
	case "npm":
		cmdName = "npm"
		target := fmt.Sprintf("%s@%s", pkgName, ver)
		args = []string{"install", "-g", target}
	case "choco":
		cmdName = "choco"
		if ver == "latest" {
			args = []string{"upgrade", "-y", pkgName}
		} else {
			args = []string{"install", "-y", pkgName, "--version", ver, "--allow-downgrade"}
		}
	case "scoop":
		cmdName = "scoop"
		if ver == "latest" {
			args = []string{"update", pkgName}
		} else {
			args = []string{"install", fmt.Sprintf("%s@%s", pkgName, ver)}
		}
	default:
		return &InstallResult{
			Success: false,
			Error:   fmt.Sprintf("unsupported package manager: %s", req.Manager),
		}
	}

	fullCmd := fmt.Sprintf("%s %s", cmdName, strings.Join(args, " "))
	cmd := exec.Command(cmdName, args...)
	hideConsole(cmd)
	out, err := cmd.CombinedOutput()

	res := &InstallResult{
		Command: fullCmd,
		Output:  string(out),
		Success: err == nil,
	}
	if err != nil {
		res.Error = err.Error()
	}
	return res
}

// ExecuteBatchCommands sequentially runs a list of installation commands (for profile import).
func ExecuteBatchCommands(commands []string) *BatchExecutionResult {
	result := &BatchExecutionResult{
		TotalSteps: len(commands),
		Steps:      make([]BatchExecutionStep, 0, len(commands)),
		Success:    true,
	}

	for _, cmdStr := range commands {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue
		}

		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		step := BatchExecutionStep{Command: cmdStr}
		cmd := exec.Command(parts[0], parts[1:]...)
		hideConsole(cmd)
		out, err := cmd.CombinedOutput()
		step.Output = string(out)
		step.Success = err == nil
		if err != nil {
			step.Error = err.Error()
			result.Success = false
		}
		result.Steps = append(result.Steps, step)
		result.CompletedSteps++
	}

	return result
}
