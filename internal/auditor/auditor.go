package auditor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"packetinstall/internal/model"
)

type EolCycle struct {
	Cycle   string `json:"cycle"`
	Latest  string `json:"latest"`
	EolDate string `json:"eol"` // string date or unmarshaled
	IsEol   bool   `json:"-"`
}

type rawEolCycle struct {
	Cycle  string          `json:"cycle"`
	Latest string          `json:"latest"`
	Eol    json.RawMessage `json:"eol"`
}

type AuditorClient struct {
	httpClient     *http.Client
	baseUrl        string
	npmRegistryUrl string
	cache          *Cache
}

type Option func(*AuditorClient)

func WithBaseUrl(url string) Option {
	return func(a *AuditorClient) {
		a.baseUrl = url
	}
}

func WithNpmRegistryUrl(url string) Option {
	return func(a *AuditorClient) {
		a.npmRegistryUrl = url
	}
}

func NewAuditorClient(opts ...Option) *AuditorClient {
	client := &AuditorClient{
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		baseUrl:        "https://endoflife.date/api",
		npmRegistryUrl: "https://registry.npmjs.org",
		cache:          NewCache(12 * time.Hour),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// CheckRuntimeEol checks endoflife.date for lifecycle status of a runtime (node, python, go).
func (a *AuditorClient) CheckRuntimeEol(product, currentVersion string) (*model.AuditItem, error) {
	cacheKey := fmt.Sprintf("eol:%s", product)
	var cycles []EolCycle

	if cached, ok := a.cache.Get(cacheKey); ok {
		cycles = cached.([]EolCycle)
	} else {
		url := fmt.Sprintf("%s/%s.json", a.baseUrl, product)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "packetinstall/1.0")

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("upstream API returned status %d", resp.StatusCode)
		}

		var rawCycles []rawEolCycle
		if err := json.NewDecoder(resp.Body).Decode(&rawCycles); err != nil {
			return nil, err
		}

		for _, rc := range rawCycles {
			c := EolCycle{
				Cycle:  rc.Cycle,
				Latest: rc.Latest,
			}
			// Parse eol field (can be bool or date string)
			var b bool
			var s string
			if err := json.Unmarshal(rc.Eol, &b); err == nil {
				c.IsEol = b
			} else if err := json.Unmarshal(rc.Eol, &s); err == nil {
				c.EolDate = s
				// Compare with today's date
				today := time.Now().Format("2006-01-02")
				c.IsEol = today > s
			}
			cycles = append(cycles, c)
		}
		a.cache.Set(cacheKey, cycles)
	}

	cleanCurrent := strings.TrimPrefix(currentVersion, "v")
	major := extractMajor(cleanCurrent)

	for _, cycle := range cycles {
		if cycle.Cycle == major {
			status := "HEALTHY"
			msg := "Runtime cycle is active."
			updateCmd := ""

			if cycle.IsEol {
				status = "EOL_CRITICAL"
				msg = fmt.Sprintf("%s %s reached End-Of-Life (EOL %s). Security risks!", product, cycle.Cycle, cycle.EolDate)
				updateCmd = fmt.Sprintf("Upgrade %s to active LTS release.", product)
			} else if isVersionNewer(cycle.Latest, cleanCurrent) {
				status = "UPDATE_AVAILABLE"
				msg = fmt.Sprintf("Latest patch in cycle %s is %s.", cycle.Cycle, cycle.Latest)
				updateCmd = fmt.Sprintf("Upgrade to %s", cycle.Latest)
			}

			return &model.AuditItem{
				Name:           product,
				Type:           "runtime",
				CurrentVersion: currentVersion,
				LatestVersion:  cycle.Latest,
				Status:         status,
				Message:        msg,
				UpdateCommand:  updateCmd,
			}, nil
		}
	}

	return &model.AuditItem{
		Name:           product,
		Type:           "runtime",
		CurrentVersion: currentVersion,
		Status:         "UNKNOWN",
		Message:        "No matching lifecycle cycle found.",
	}, nil
}

// CheckNpmPackage checks latest release version on NPM registry.
func (a *AuditorClient) CheckNpmPackage(pkgName, currentVersion string) (*model.AuditItem, error) {
	cacheKey := fmt.Sprintf("npm:%s", pkgName)
	var latestVersion string

	if cached, ok := a.cache.Get(cacheKey); ok {
		latestVersion = cached.(string)
	} else {
		url := fmt.Sprintf("%s/%s/latest", a.npmRegistryUrl, pkgName)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "packetinstall/1.0")

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("npm registry returned %d", resp.StatusCode)
		}

		var meta struct {
			Version string `json:"version"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
			return nil, err
		}
		latestVersion = meta.Version
		a.cache.Set(cacheKey, latestVersion)
	}

	cleanCurrent := strings.TrimPrefix(currentVersion, "v")
	cleanLatest := strings.TrimPrefix(latestVersion, "v")

	status := "HEALTHY"
	msg := "Package is up to date."
	updateCmd := ""

	if isVersionNewer(cleanLatest, cleanCurrent) {
		status = "UPDATE_AVAILABLE"
		msg = fmt.Sprintf("New version available: %s (current: %s)", cleanLatest, cleanCurrent)
		updateCmd = fmt.Sprintf("npm install -g %s@latest", pkgName)
	}

	return &model.AuditItem{
		Name:           pkgName,
		Type:           "package",
		Manager:        "npm",
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		Status:         status,
		Message:        msg,
		UpdateCommand:  updateCmd,
	}, nil
}

func extractMajor(v string) string {
	parts := strings.Split(v, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return v
}

func isVersionNewer(latest, current string) bool {
	lParts := strings.Split(latest, ".")
	cParts := strings.Split(current, ".")

	maxLen := len(lParts)
	if len(cParts) > maxLen {
		maxLen = len(cParts)
	}

	for i := range maxLen {
		lNum := 0
		cNum := 0
		if i < len(lParts) {
			lNum, _ = strconv.Atoi(lParts[i])
		}
		if i < len(cParts) {
			cNum, _ = strconv.Atoi(cParts[i])
		}
		if lNum > cNum {
			return true
		}
		if lNum < cNum {
			return false
		}
	}
	return false
}
