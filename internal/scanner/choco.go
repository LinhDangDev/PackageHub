package scanner

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"packetinstall/internal/model"
)

type nuspecMetadata struct {
	ID          string `xml:"id"`
	Version     string `xml:"version"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

type nuspecPackage struct {
	Metadata nuspecMetadata `xml:"metadata"`
}

// ScanChoco parses Chocolatey installed package folders and their .nuspec files.
func ScanChoco(baseDir string) ([]model.Package, error) {
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
		pkgName := entry.Name()
		nuspecPath := filepath.Join(baseDir, pkgName, pkgName+".nuspec")
		data, err := os.ReadFile(nuspecPath)
		if err != nil {
			// Fallback: search for any .nuspec in the folder
			files, _ := os.ReadDir(filepath.Join(baseDir, pkgName))
			for _, f := range files {
				if filepath.Ext(f.Name()) == ".nuspec" {
					nuspecPath = filepath.Join(baseDir, pkgName, f.Name())
					data, err = os.ReadFile(nuspecPath)
					break
				}
			}
		}

		if err == nil {
			var spec nuspecPackage
			if err := xml.Unmarshal(data, &spec); err == nil && spec.Metadata.ID != "" {
				packages = append(packages, model.Package{
					Manager:     "choco",
					Name:        spec.Metadata.ID,
					Version:     spec.Metadata.Version,
					Description: spec.Metadata.Description,
					Path:        filepath.Join(baseDir, pkgName),
				})
			}
		}
	}

	return packages, nil
}
