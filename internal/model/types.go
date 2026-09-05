package model

import "time"

// Package represents an installed software package, runtime, or CLI tool.
type Package struct {
	Manager     string              `json:"manager"`     // choco, scoop, npm, pipx, cargo, winget, runtime
	Name        string              `json:"name"`        // e.g. ripgrep, @anthropic-ai/claude-code
	Version     string              `json:"version"`     // e.g. 14.1.0
	Description string              `json:"description"` // brief description
	Path        string              `json:"path"`        // install path or manifest location
	Diagnostics *PackageDiagnostics `json:"diagnostics,omitempty"`
}

// PackageIssue represents a specific problem or warning detected with a package.
type PackageIssue struct {
	Severity string `json:"severity"` // "ERROR", "WARNING", "INFO"
	Code     string `json:"code"`     // "PATH_MISSING", "OUTDATED", "ENGINE_MISMATCH", "EOL"
	Message  string `json:"message"`
	FixHint  string `json:"fix_hint"`
}

// PackageDiagnostics holds real-time health, update status, and detected issues.
type PackageDiagnostics struct {
	Status        string         `json:"status"` // "HEALTHY", "OUTDATED", "EOL_CRITICAL", "BROKEN", "UNKNOWN"
	LatestVersion string         `json:"latest_version"`
	UpdateCommand string         `json:"update_command"`
	InPath        bool           `json:"in_path"`
	BinaryName    string         `json:"binary_name"`
	EngineReq     string         `json:"engine_req"`
	Issues        []PackageIssue `json:"issues"`
	CheckedAt     time.Time      `json:"checked_at"`
}

// ProjectDependency represents a single library/dependency inside a project.
type ProjectDependency struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	LatestVer string `json:"latest_ver,omitempty"`
	IsDev     bool   `json:"is_dev"`
	Status    string `json:"status"` // "OK", "OUTDATED", "DEPRECATED", "RISKY"
	Issue     string `json:"issue,omitempty"`
}

// ProjectInfo represents a code repository or project found on disk.
type ProjectInfo struct {
	Name          string              `json:"name"`
	Path          string              `json:"path"`
	Language      string              `json:"language"`     // "TypeScript", "JavaScript", "Go", "Python", "Rust"
	Framework     string              `json:"framework"`    // "Next.js", "React", "Node.js", "Gin", "FastAPI", etc.
	HealthScore   int                 `json:"health_score"` // 0 - 100
	Dependencies  []ProjectDependency `json:"dependencies"`
	OutdatedCount int                 `json:"outdated_count"`
	IssueCount    int                 `json:"issue_count"`
}

// ProjectScanResult represents the result of scanning a drive or root directory.
type ProjectScanResult struct {
	ScanPath   string        `json:"scan_path"`
	Projects   []ProjectInfo `json:"projects"`
	TotalDeps  int           `json:"total_deps"`
	DurationMs int64         `json:"duration_ms"`
}

// InstallPackageRequest represents a request to install a package via CLI.
type InstallPackageRequest struct {
	Manager string `json:"manager"` // "npm", "choco", "scoop"
	Package string `json:"package"` // e.g. "lodash" or "BurntSushi.ripgrep.MSVC"
	Global  bool   `json:"global"`
}

// SkillFlag represents a command-line flag or option for a skill.
type SkillFlag struct {
	Name        string `json:"name"`             // e.g. "--mode", "--fast", "--tdd"
	Values      string `json:"values,omitempty"` // e.g. "search | creative | wild | all"
	Description string `json:"description"`      // e.g. "Bỏ qua phỏng vấn kiểm tra phong cách ban đầu"
}

// Skill represents an AI Agent skill or workflow extension.
type Skill struct {
	Name         string      `json:"name"`          // e.g. ck:ask, sequential-thinking
	ToolOrigin   string      `json:"tool_origin"`   // "claudekit", "opencode", "codex", "universal", "local"
	Category     string      `json:"category"`      // "utilities", "planning", "coding", "design", etc.
	Tier         string      `json:"tier"`          // "beginner", "intermediate", "pro"
	Command      string      `json:"command"`       // "/ck:ask" or "/ask"
	ArgumentHint string      `json:"argument_hint"` // "[technical-question]"
	WhenToUse    string      `json:"when_to_use"`
	Path         string      `json:"path"`          // filesystem path
	Description  string      `json:"description"`   // extracted from SKILL.md
	GitRemote    string      `json:"git_remote"`    // git origin URL if git-backed
	CommitSHA    string      `json:"commit_sha"`    // HEAD commit hash if git-backed
	FullDocs     string      `json:"full_docs"`     // Markdown content for VividKit-style docs
	Examples     []string    `json:"examples"`      // Sample prompt commands
	Flags        []SkillFlag `json:"flags"`         // Parsed flags and their meanings
}

// McpServer represents a Model Context Protocol server configuration.
type McpServer struct {
	Name     string            `json:"name"`     // e.g. filesystem, github
	Source   string            `json:"source"`   // claude-desktop, cursor, windsurf
	Command  string            `json:"command"`  // e.g. npx, python, docker
	Args     []string          `json:"args"`     // execution arguments
	Env      map[string]string `json:"env"`      // environment variables
	Disabled bool              `json:"disabled"` // whether server is temporarily disabled
}

// SystemState represents the unified snapshot of everything scanned on the machine.
type SystemState struct {
	Packages   []Package   `json:"packages"`
	Skills     []Skill     `json:"skills"`
	McpServers []McpServer `json:"mcp_servers"`
	ScannedAt  time.Time   `json:"scanned_at"`
	DurationMs int64       `json:"duration_ms"`
}

// AuditItem represents the version health and update assessment for a tool.
type AuditItem struct {
	Name           string `json:"name"`
	Type           string `json:"type"` // runtime, package, skill
	Manager        string `json:"manager"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Status         string `json:"status"` // HEALTHY, UPDATE_AVAILABLE, EOL_CRITICAL, UNKNOWN
	Message        string `json:"message"`
	UpdateCommand  string `json:"update_command"`
}

// Profile represents the declarative exportable/importable machine configuration.
type Profile struct {
	SchemaVersion string                 `yaml:"schema_version" json:"schema_version"`
	Metadata      ProfileMetadata        `yaml:"metadata" json:"metadata"`
	Runtimes      []RuntimeSpec          `yaml:"runtimes,omitempty" json:"runtimes,omitempty"`
	SystemTools   []SystemPackageSpec    `yaml:"system_packages,omitempty" json:"system_packages,omitempty"`
	GlobalCLIs    map[string][]string    `yaml:"global_clis,omitempty" json:"global_clis,omitempty"`
	McpServers    []McpServerProfileSpec `yaml:"mcp_servers,omitempty" json:"mcp_servers,omitempty"`
	Skills        []SkillProfileSpec     `yaml:"skills,omitempty" json:"skills,omitempty"`
}

type ProfileMetadata struct {
	Name        string `yaml:"name" json:"name"`
	CreatedAt   string `yaml:"created_at" json:"created_at"`
	SourceOS    string `yaml:"source_os" json:"source_os"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type RuntimeSpec struct {
	ID                string `yaml:"id" json:"id"`
	VersionConstraint string `yaml:"version_constraint" json:"version_constraint"`
	PreferredManager  string `yaml:"preferred_manager,omitempty" json:"preferred_manager,omitempty"`
}

type SystemPackageSpec struct {
	ID      string                    `yaml:"id" json:"id"`
	Windows *PlatformPackageSpec      `yaml:"windows,omitempty" json:"windows,omitempty"`
	MacOS   *PlatformPackageSpec      `yaml:"macos,omitempty" json:"macos,omitempty"`
	Linux   *PlatformPackageSpec      `yaml:"linux,omitempty" json:"linux,omitempty"`
}

type PlatformPackageSpec struct {
	Manager   string `yaml:"manager" json:"manager"` // winget, scoop, choco, brew, apt
	PackageID string `yaml:"package_id" json:"package_id"`
}

type McpServerProfileSpec struct {
	Name    string            `yaml:"name" json:"name"`
	Command string            `yaml:"command" json:"command"`
	Args    []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Env     map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
}

type SkillProfileSpec struct {
	Name      string `yaml:"name" json:"name"`
	RepoURL   string `yaml:"repo_url,omitempty" json:"repo_url,omitempty"`
	CommitSHA string `yaml:"commit_sha,omitempty" json:"commit_sha,omitempty"`
}
