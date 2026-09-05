package scanner

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"packetinstall/internal/model"
)

type rawPackageJson struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

var knownDeprecatedPackages = map[string]string{
	"request":       "Package 'request' has been deprecated since 2020.",
	"moment":        "Package 'moment' is in maintenance mode. Consider luxon or date-fns.",
	"querystring":   "Built-in querystring is legacy. Use URLSearchParams.",
	"colors":        "Package 'colors' suffered supply chain disruption.",
	"left-pad":      "Obsolete single-line package.",
	"core-js@2":     "core-js@2 is deprecated.",
	"faker":         "Unmaintained legacy faker.",
}

// ScanProjects walks the directory tree up to maxDepth and audits all discovered code projects.
func ScanProjects(rootPath string, maxDepth int) (*model.ProjectScanResult, error) {
	start := time.Now()
	result := &model.ProjectScanResult{
		ScanPath: rootPath,
		Projects: make([]model.ProjectInfo, 0),
	}

	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return result, nil
	}

	visitedDirs := make(map[string]bool)

	var walkFn func(currentPath string, depth int)
	walkFn = func(currentPath string, depth int) {
		if depth > maxDepth {
			return
		}

		baseName := filepath.Base(currentPath)
		if isIgnoredDir(baseName) {
			return
		}

		// Check if current directory is a project root
		proj, isProj := inspectProjectDir(currentPath)
		if isProj {
			result.Projects = append(result.Projects, *proj)
			result.TotalDeps += len(proj.Dependencies)
			visitedDirs[currentPath] = true
			// Don't recurse deeply into recognized non-monorepo project subfolders
			if !hasSubprojects(currentPath) {
				return
			}
		}

		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				walkFn(filepath.Join(currentPath, entry.Name()), depth+1)
			}
		}
	}

	walkFn(rootPath, 0)
	result.DurationMs = time.Since(start).Milliseconds()

	return result, nil
}

func inspectProjectDir(dir string) (*model.ProjectInfo, bool) {
	pkgJsonPath := filepath.Join(dir, "package.json")
	goModPath := filepath.Join(dir, "go.mod")
	cargoPath := filepath.Join(dir, "Cargo.toml")
	pyPath := filepath.Join(dir, "requirements.txt")

	if _, err := os.Stat(pkgJsonPath); err == nil {
		return parseNodeProject(dir, pkgJsonPath)
	}

	if _, err := os.Stat(goModPath); err == nil {
		return parseGoProject(dir, goModPath)
	}

	if _, err := os.Stat(cargoPath); err == nil {
		return parseCargoProject(dir, cargoPath)
	}

	if _, err := os.Stat(pyPath); err == nil {
		return parsePythonProject(dir, pyPath)
	}

	return nil, false
}

func parseNodeProject(dir, pkgJsonPath string) (*model.ProjectInfo, bool) {
	data, err := os.ReadFile(pkgJsonPath)
	if err != nil {
		return nil, false
	}

	var pj rawPackageJson
	if err := json.Unmarshal(data, &pj); err != nil {
		return nil, false
	}

	projName := pj.Name
	if projName == "" {
		projName = filepath.Base(dir)
	}

	// Detect language: check for tsconfig.json or typescript in deps
	lang := "JavaScript"
	if _, err := os.Stat(filepath.Join(dir, "tsconfig.json")); err == nil {
		lang = "TypeScript"
	} else if _, ok := pj.DevDependencies["typescript"]; ok {
		lang = "TypeScript"
	}

	// Detect framework
	framework := "Node.js"
	allDeps := make(map[string]string)
	for k, v := range pj.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pj.DevDependencies {
		allDeps[k] = v
	}

	if _, ok := allDeps["next"]; ok {
		framework = "Next.js"
	} else if _, ok := allDeps["react"]; ok {
		framework = "React"
	} else if _, ok := allDeps["vue"]; ok {
		framework = "Vue"
	} else if _, ok := allDeps["svelte"]; ok {
		framework = "Svelte"
	} else if _, ok := allDeps["express"]; ok {
		framework = "Express"
	} else if _, ok := allDeps["turbo"]; ok {
		framework = "Turborepo"
	}

	deps := make([]model.ProjectDependency, 0, len(allDeps))
	outdatedCount := 0
	issueCount := 0

	for name, ver := range pj.Dependencies {
		dep := auditDependency(name, ver, false)
		if dep.Status != "OK" {
			issueCount++
			if dep.Status == "OUTDATED" {
				outdatedCount++
			}
		}
		deps = append(deps, dep)
	}

	for name, ver := range pj.DevDependencies {
		dep := auditDependency(name, ver, true)
		if dep.Status != "OK" {
			issueCount++
		}
		deps = append(deps, dep)
	}

	healthScore := calculateHealthScore(issueCount, len(deps))

	return &model.ProjectInfo{
		Name:          projName,
		Path:          dir,
		Language:      lang,
		Framework:     framework,
		HealthScore:   healthScore,
		Dependencies:  deps,
		OutdatedCount: outdatedCount,
		IssueCount:    issueCount,
	}, true
}

func parseGoProject(dir, goModPath string) (*model.ProjectInfo, bool) {
	file, err := os.Open(goModPath)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	moduleName := filepath.Base(dir)
	framework := "Go Standard"
	deps := make([]model.ProjectDependency, 0)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		} else if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		} else if inRequire && line == ")" {
			inRequire = false
			continue
		}

		if inRequire || strings.HasPrefix(line, "require ") {
			cleanLine := strings.TrimPrefix(line, "require ")
			parts := strings.Fields(cleanLine)
			if len(parts) >= 2 {
				depName := parts[0]
				depVer := parts[1]
				if strings.Contains(depName, "gin-gonic/gin") {
					framework = "Gin"
				} else if strings.Contains(depName, "fiber") {
					framework = "Fiber"
				} else if strings.Contains(depName, "webview") {
					framework = "WebView Desktop"
				}
				deps = append(deps, model.ProjectDependency{
					Name:    depName,
					Version: depVer,
					Status:  "OK",
				})
			}
		}
	}

	return &model.ProjectInfo{
		Name:         moduleName,
		Path:         dir,
		Language:     "Go",
		Framework:    framework,
		HealthScore:  100,
		Dependencies: deps,
	}, true
}

func parseCargoProject(dir, cargoPath string) (*model.ProjectInfo, bool) {
	return &model.ProjectInfo{
		Name:        filepath.Base(dir),
		Path:        dir,
		Language:    "Rust",
		Framework:   "Cargo Crate",
		HealthScore: 100,
	}, true
}

func parsePythonProject(dir, pyPath string) (*model.ProjectInfo, bool) {
	data, err := os.ReadFile(pyPath)
	if err != nil {
		return nil, false
	}
	lines := strings.Split(string(data), "\n")
	deps := make([]model.ProjectDependency, 0)
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			parts := strings.Split(trimmed, "==")
			name := parts[0]
			ver := "latest"
			if len(parts) > 1 {
				ver = parts[1]
			}
			deps = append(deps, model.ProjectDependency{
				Name:    name,
				Version: ver,
				Status:  "OK",
			})
		}
	}
	return &model.ProjectInfo{
		Name:         filepath.Base(dir),
		Path:         dir,
		Language:     "Python",
		Framework:    "Python Package",
		HealthScore:  100,
		Dependencies: deps,
	}, true
}

func auditDependency(name, ver string, isDev bool) model.ProjectDependency {
	status := "OK"
	issue := ""

	if msg, deprecated := knownDeprecatedPackages[name]; deprecated {
		status = "DEPRECATED"
		issue = msg
	} else if ver == "*" || ver == "latest" {
		status = "RISKY"
		issue = "Wildcard version can introduce breaking changes unexpectedly."
	}

	return model.ProjectDependency{
		Name:    name,
		Version: ver,
		IsDev:   isDev,
		Status:  status,
		Issue:   issue,
	}
}

func calculateHealthScore(issueCount, totalDeps int) int {
	if totalDeps == 0 {
		return 100
	}
	penalty := issueCount * 12
	score := 100 - penalty
	if score < 20 {
		score = 20
	}
	return score
}

func isIgnoredDir(name string) bool {
	ignored := map[string]bool{
		"node_modules": true,
		".git":         true,
		".turbo":       true,
		".vscode":      true,
		".idea":        true,
		"dist":         true,
		"build":        true,
		"target":       true,
		"vendor":       true,
		".venv":        true,
		"__pycache__":  true,
		".next":        true,
		".cache":       true,
	}
	return ignored[name]
}

func hasSubprojects(dir string) bool {
	// Turborepo / pnpm-workspace / packages dir
	if _, err := os.Stat(filepath.Join(dir, "pnpm-workspace.yaml")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(dir, "turbo.json")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(dir, "packages")); err == nil {
		return true
	}
	return false
}
