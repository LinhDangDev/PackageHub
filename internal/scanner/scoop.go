package scanner

import (
	"os"
	"path/filepath"

	"packetinstall/internal/model"
)

// ScanScoop parses Scoop apps installed in the scoop apps directory.
func ScanScoop(baseDir string) ([]model.Package, error) {
	var packages []model.Package
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return packages, nil
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "scoop" {
			continue
		}
		appName := entry.Name()
		appPath := filepath.Join(baseDir, appName)

		// Scoop links the active version to "current"
		currentLink := filepath.Join(appPath, "current")
		version := "unknown"
		if target, err := os.Readlink(currentLink); err == nil {
			version = filepath.Base(target)
		} else {
			// Fallback: pick the latest subfolder that is not "current"
			subEntries, _ := os.ReadDir(appPath)
			for _, sub := range subEntries {
				if sub.IsDir() && sub.Name() != "current" {
					version = sub.Name()
				}
			}
		}

		packages = append(packages, model.Package{
			Manager: "scoop",
			Name:    appName,
			Version: version,
			Path:    appPath,
		})
	}

	return packages, nil
}
