package cleaner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"packetinstall/internal/model"
)

// hideConsole prevents terminal popups on Windows.
func hideConsole(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
	}
}

// calcDirSize recursively computes the total byte size of a directory.
func calcDirSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// ================= 1. GEEK-STYLE LEFTOVERS SCAVENGER =================

// normalizeToolName strips scopes and slashes for filesystem/registry queries.
func normalizeToolName(name string) string {
	name = strings.TrimSpace(name)
	if idx := strings.LastIndex(name, "/"); idx != -1 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, "\\"); idx != -1 {
		name = name[idx+1:]
	}
	return strings.ToLower(name)
}

// ScanLeftovers searches for leftover directories, registry traces, and dead PATH entries.
func ScanLeftovers(toolName string) model.LeftoverReport {
	norm := normalizeToolName(toolName)
	report := model.LeftoverReport{
		ToolName: toolName,
		Items:    []model.LeftoverItem{},
	}
	if norm == "" {
		return report
	}

	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")
	programData := os.Getenv("ProgramData")
	userProfile := os.Getenv("USERPROFILE")

	// 1. Filesystem candidate paths
	candidates := []struct {
		path string
		desc string
	}{
		{filepath.Join(appData, norm), "Roaming configuration & cache"},
		{filepath.Join(localAppData, norm), "Local application data & logs"},
		{filepath.Join(localAppData, "Programs", norm), "User-level binary directory"},
		{filepath.Join(programData, norm), "System-wide program data"},
		{filepath.Join(userProfile, "."+norm), "User home dotfile/dotfolder"},
		{filepath.Join(userProfile, norm), "User home application directory"},
		{filepath.Join(localAppData, "Temp", norm), "Temporary execution files"},
	}

	idCounter := 1
	for _, c := range candidates {
		if c.path == "" {
			continue
		}
		if fi, err := os.Stat(c.path); err == nil {
			itemType := "file"
			var size int64
			if fi.IsDir() {
				itemType = "dir"
				size = calcDirSize(c.path)
			} else {
				size = fi.Size()
			}
			report.Items = append(report.Items, model.LeftoverItem{
				ID:          fmt.Sprintf("item-%d", idCounter),
				Type:        itemType,
				Path:        c.path,
				Description: c.desc,
				Size:        size,
				Selected:    true,
			})
			report.TotalSize += size
			idCounter++
		}
	}

	// 2. Registry candidate keys (Windows)
	if runtime.GOOS == "windows" {
		regPaths := []string{
			`HKCU\Software\` + norm,
			`HKLM\Software\` + norm,
			`HKCU\Software\` + strings.Title(norm),
			`HKLM\Software\` + strings.Title(norm),
		}

		for _, regPath := range regPaths {
			cmd := exec.Command("reg", "query", regPath)
			hideConsole(cmd)
			if err := cmd.Run(); err == nil {
				report.Items = append(report.Items, model.LeftoverItem{
					ID:          fmt.Sprintf("item-%d", idCounter),
					Type:        "registry",
					Path:        regPath,
					Description: "Windows Registry configuration key",
					Size:        0,
					Selected:    true,
				})
				idCounter++
			}
		}
	}

	// 3. Leftover / Dangling PATH entries
	pathReport := AuditPath()
	for _, entry := range pathReport.Entries {
		lower := strings.ToLower(entry.Path)
		if strings.Contains(lower, norm) || !entry.Exists {
			if strings.Contains(lower, norm) {
				report.Items = append(report.Items, model.LeftoverItem{
					ID:          fmt.Sprintf("item-%d", idCounter),
					Type:        "path",
					Path:        entry.Path,
					Description: "Dangling or related PATH environment entry",
					Size:        0,
					Selected:    true,
				})
				idCounter++
			}
		}
	}

	report.TotalItems = len(report.Items)
	return report
}

// PurgeLeftovers permanently deletes selected leftover files, registry keys, and path entries.
func PurgeLeftovers(req model.PurgeLeftoversRequest) model.PurgeResult {
	res := model.PurgeResult{Success: true}
	fullReport := ScanLeftovers(req.ToolName)

	targetsByID := make(map[string]model.LeftoverItem)
	for _, it := range fullReport.Items {
		targetsByID[it.ID] = it
	}

	// If no specific IDs provided, purge all detected
	var toPurge []model.LeftoverItem
	if len(req.ItemIDs) == 0 {
		toPurge = fullReport.Items
	} else {
		for _, id := range req.ItemIDs {
			if it, ok := targetsByID[id]; ok {
				toPurge = append(toPurge, it)
			}
		}
	}

	for _, item := range toPurge {
		switch item.Type {
		case "dir", "file":
			if err := os.RemoveAll(item.Path); err != nil {
				res.Errors = append(res.Errors, fmt.Sprintf("Failed to delete %s: %v", item.Path, err))
			} else {
				res.PurgedCount++
			}
		case "registry":
			cmd := exec.Command("reg", "delete", item.Path, "/f")
			hideConsole(cmd)
			if err := cmd.Run(); err != nil {
				res.Errors = append(res.Errors, fmt.Sprintf("Failed to delete registry %s: %v", item.Path, err))
			} else {
				res.PurgedCount++
			}
		case "path":
			_ = removeEntryFromUserPath(item.Path)
			res.PurgedCount++
		}
	}

	if len(res.Errors) > 0 && res.PurgedCount == 0 {
		res.Success = false
	}
	return res
}

// ================= 2. ZOMBIE PATH PURGER =================

// AuditPath inspects all user and system PATH entries to detect dead/missing directories.
func AuditPath() model.PathAuditReport {
	report := model.PathAuditReport{
		Entries: []model.PathEntry{},
	}

	rawPath := os.Getenv("PATH")
	entries := strings.Split(rawPath, string(os.PathListSeparator))

	seen := make(map[string]bool)
	for _, p := range entries {
		p = strings.TrimSpace(p)
		if p == "" || seen[strings.ToLower(p)] {
			continue
		}
		seen[strings.ToLower(p)] = true

		exists := true
		var errStr string
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			exists = false
			if err != nil {
				errStr = err.Error()
			} else {
				errStr = "not a directory"
			}
			report.DeadCount++
		}

		report.Entries = append(report.Entries, model.PathEntry{
			Path:   p,
			Scope:  "user",
			Exists: exists,
			Error:  errStr,
		})
	}

	report.TotalEntries = len(report.Entries)
	return report
}

// PruneDeadPaths cleans invalid directories from the Windows User PATH environment variable.
func PruneDeadPaths() (model.PathAuditReport, error) {
	if runtime.GOOS != "windows" {
		return AuditPath(), nil
	}

	// Fetch current User PATH directly via PowerShell
	cmdGet := exec.Command("powershell", "-NoProfile", "-Command", "[Environment]::GetEnvironmentVariable('Path', 'User')")
	hideConsole(cmdGet)
	out, err := cmdGet.Output()
	if err != nil {
		return AuditPath(), err
	}

	rawUserPath := strings.TrimSpace(string(out))
	parts := strings.Split(rawUserPath, ";")

	var validParts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Expand environment variables like %USERPROFILE% if present
		expanded := os.ExpandEnv(p)
		if fi, err := os.Stat(expanded); err == nil && fi.IsDir() {
			validParts = append(validParts, p)
		}
	}

	cleanedPath := strings.Join(validParts, ";")
	psSet := fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', '%s', 'User')", strings.ReplaceAll(cleanedPath, "'", "''"))
	cmdSet := exec.Command("powershell", "-NoProfile", "-Command", psSet)
	hideConsole(cmdSet)
	_ = cmdSet.Run()

	return AuditPath(), nil
}

func removeEntryFromUserPath(target string) error {
	if runtime.GOOS != "windows" {
		return nil
	}
	cmdGet := exec.Command("powershell", "-NoProfile", "-Command", "[Environment]::GetEnvironmentVariable('Path', 'User')")
	hideConsole(cmdGet)
	out, err := cmdGet.Output()
	if err != nil {
		return err
	}

	rawUserPath := strings.TrimSpace(string(out))
	parts := strings.Split(rawUserPath, ";")
	var surviving []string
	targetLower := strings.ToLower(strings.TrimSpace(target))

	for _, p := range parts {
		pTrim := strings.TrimSpace(p)
		if strings.ToLower(pTrim) != targetLower && pTrim != "" {
			surviving = append(surviving, pTrim)
		}
	}

	newPath := strings.Join(surviving, ";")
	psSet := fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', '%s', 'User')", strings.ReplaceAll(newPath, "'", "''"))
	cmdSet := exec.Command("powershell", "-NoProfile", "-Command", psSet)
	hideConsole(cmdSet)
	return cmdSet.Run()
}

// ================= 3. DEV CACHE & STORAGE HOG CLEANER =================

// AuditDevCaches inspects disk usage across developer caches (NPM, Pip, Go, Cargo, etc.).
func AuditDevCaches() model.DevCacheReport {
	report := model.DevCacheReport{
		Caches: []model.DevCacheItem{},
	}

	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")
	userProfile := os.Getenv("USERPROFILE")

	cacheSpecs := []struct {
		id       string
		name     string
		path     string
		cleanCmd string
	}{
		{"npm", "NPM Cache", filepath.Join(appData, "npm-cache"), "npm cache clean --force"},
		{"pip", "Python Pip Cache", filepath.Join(localAppData, "pip", "cache"), "pip cache purge"},
		{"go_build", "Go Build Cache", filepath.Join(localAppData, "go-build"), "go clean -cache"},
		{"go_mod", "Go Module Cache", filepath.Join(userProfile, "go", "pkg", "mod", "cache"), "go clean -modcache"},
		{"cargo", "Rust Cargo Cache", filepath.Join(userProfile, ".cargo", "registry", "cache"), "cargo clean"},
		{"yarn", "Yarn Cache", filepath.Join(localAppData, "Yarn", "Cache"), "yarn cache clean"},
		{"pnpm", "pnpm Store", filepath.Join(localAppData, "pnpm", "store"), "pnpm store prune"},
		{"choco", "Chocolatey Cache", filepath.Join(localAppData, "Chocolatey", "cache"), "rm"},
	}

	for _, spec := range cacheSpecs {
		if spec.path == "" {
			continue
		}
		item := model.DevCacheItem{
			ID:       spec.id,
			Name:     spec.name,
			Path:     spec.path,
			CleanCmd: spec.cleanCmd,
		}

		if fi, err := os.Stat(spec.path); err == nil && fi.IsDir() {
			item.Exists = true
			item.Size = calcDirSize(spec.path)
			report.TotalSize += item.Size
		}

		report.Caches = append(report.Caches, item)
	}

	return report
}

// CleanDevCache removes or cleans a specific developer cache and returns bytes freed.
func CleanDevCache(cacheID string) (int64, error) {
	report := AuditDevCaches()
	var target *model.DevCacheItem
	for i := range report.Caches {
		if report.Caches[i].ID == cacheID {
			target = &report.Caches[i]
			break
		}
	}

	if target == nil || !target.Exists {
		return 0, fmt.Errorf("cache '%s' does not exist or not found", cacheID)
	}

	freed := target.Size

	// Try dedicated command first if available
	cleanedViaCmd := false
	if target.CleanCmd != "rm" && target.CleanCmd != "" {
		parts := strings.Fields(target.CleanCmd)
		if len(parts) > 0 {
			cmd := exec.Command(parts[0], parts[1:]...)
			hideConsole(cmd)
			if err := cmd.Run(); err == nil {
				cleanedViaCmd = true
			}
		}
	}

	// Fallback to directory contents removal
	if !cleanedViaCmd {
		entries, err := os.ReadDir(target.Path)
		if err == nil {
			for _, e := range entries {
				_ = os.RemoveAll(filepath.Join(target.Path, e.Name()))
			}
		}
	}

	return freed, nil
}

// ================= 4. DEV PORT & PROCESS CONFLICT AUDITOR =================

// AuditDevPorts lists active listening TCP ports, matching them to process names and PIDs.
func AuditDevPorts() model.DevPortReport {
	report := model.DevPortReport{
		Ports: []model.DevPortItem{},
	}

	if runtime.GOOS != "windows" {
		return report
	}

	// 1. Build PID -> Process Name map in a single fast call (< 30ms)
	pidToName := make(map[int]string)
	cmdTaskList := exec.Command("tasklist", "/fo", "csv", "/nh")
	hideConsole(cmdTaskList)
	if out, err := cmdTaskList.Output(); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, "\",\"")
			if len(parts) >= 2 {
				name := strings.Trim(parts[0], "\"")
				pidStr := strings.Trim(parts[1], "\"")
				if pid, err := strconv.Atoi(pidStr); err == nil {
					pidToName[pid] = name
				}
			}
		}
	}

	// 2. Query listening TCP ports via netstat
	cmdNetstat := exec.Command("netstat", "-ano", "-p", "tcp")
	hideConsole(cmdNetstat)
	out, err := cmdNetstat.Output()
	if err != nil {
		return report
	}

	// Line format: TCP    0.0.0.0:3000           0.0.0.0:0              LISTENING       14200
	re := regexp.MustCompile(`^\s*TCP\s+\S+:(\d+)\s+\S+\s+LISTENING\s+(\d+)`)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	seenPorts := make(map[int]bool)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			port, _ := strconv.Atoi(matches[1])
			pid, _ := strconv.Atoi(matches[2])

			if port > 0 && !seenPorts[port] {
				seenPorts[port] = true
				procName := pidToName[pid]
				if procName == "" {
					procName = "System / Unknown"
				}

				report.Ports = append(report.Ports, model.DevPortItem{
					Port:        port,
					PID:         pid,
					ProcessName: procName,
					Protocol:    "TCP",
				})
			}
		}
	}

	report.TotalListening = len(report.Ports)
	return report
}

// KillProcess forcefully terminates a process by PID.
func KillProcess(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid PID %d", pid)
	}
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
		hideConsole(cmd)
		return cmd.Run()
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
