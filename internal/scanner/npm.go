package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"packetinstall/internal/model"
)

type npmPackageJson struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// ScanNpm parses global NPM node_modules directly without running the npm CLI.
func ScanNpm(baseDir string) ([]model.Package, error) {
	var packages []model.Package
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return packages, nil
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()

		// Check for scoped packages (@scope/package)
		if strings.HasPrefix(name, "@") {
			scopedPath := filepath.Join(baseDir, name)
			scopedEntries, err := os.ReadDir(scopedPath)
			if err == nil {
				for _, scopedEntry := range scopedEntries {
					if !scopedEntry.IsDir() {
						continue
					}
					pkgJsonPath := filepath.Join(scopedPath, scopedEntry.Name(), "package.json")
					if pkg, err := readNpmPackageJson(pkgJsonPath); err == nil {
						pkg.Manager = "npm"
						pkg.Path = filepath.Join(scopedPath, scopedEntry.Name())
						packages = append(packages, *pkg)
					}
				}
			}
		} else {
			pkgJsonPath := filepath.Join(baseDir, name, "package.json")
			if pkg, err := readNpmPackageJson(pkgJsonPath); err == nil {
				pkg.Manager = "npm"
				pkg.Path = filepath.Join(baseDir, name)
				packages = append(packages, *pkg)
			}
		}
	}

	return packages, nil
}

func readNpmPackageJson(path string) (*model.Package, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pj npmPackageJson
	if err := json.Unmarshal(data, &pj); err != nil {
		return nil, err
	}
	return &model.Package{
		Name:        pj.Name,
		Version:     pj.Version,
		Description: pj.Description,
	}, nil
}
