package web

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"runtime"

	"packetinstall/internal/auditor"
	"packetinstall/internal/installer"
	"packetinstall/internal/model"
	"packetinstall/internal/profile"
	"packetinstall/internal/scanner"
)

//go:embed static/*
var staticFS embed.FS

type Server struct {
	scannerOpts scanner.ScanOptions
	auditor     *auditor.AuditorClient
	mux         *http.ServeMux
}

func NewServer(opts scanner.ScanOptions) *Server {
	s := &Server{
		scannerOpts: opts,
		auditor:     auditor.NewAuditorClient(),
		mux:         http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /api/scan", s.handleScan)
	s.mux.HandleFunc("GET /api/audit", s.handleAudit)
	s.mux.HandleFunc("GET /api/adapters", s.handleAdapters)
	s.mux.HandleFunc("GET /api/history", s.handleHistory)
	s.mux.HandleFunc("GET /api/updates", s.handleUpdates)
	s.mux.HandleFunc("POST /api/package/diagnose", s.handlePackageDiagnose)
	s.mux.HandleFunc("POST /api/projects/scan", s.handleProjectsScan)
	s.mux.HandleFunc("POST /api/packages/install", s.handlePackageInstall)
	s.mux.HandleFunc("POST /api/packages/uninstall", s.handlePackageUninstall)
	s.mux.HandleFunc("POST /api/packages/switch-version", s.handlePackageSwitchVersion)
	s.mux.HandleFunc("POST /api/profile/execute-batch", s.handleProfileExecuteBatch)
	s.mux.HandleFunc("POST /api/profile/export-zip", s.handleProfileExportZip)
	s.mux.HandleFunc("POST /api/profile/export-selective", s.handleProfileExportSelective)
	s.mux.HandleFunc("POST /api/fix", s.handleFixIssue)
	s.mux.HandleFunc("POST /api/project/fix-dep", s.handleProjectFixDep)
	s.mux.HandleFunc("GET /api/profile/export", s.handleProfileExport)
	s.mux.HandleFunc("POST /api/profile/diff", s.handleProfileDiff)
	s.mux.HandleFunc("GET /", s.handleStatic)
}

func (s *Server) handleProjectFixDep(w http.ResponseWriter, r *http.Request) {
	var req installer.FixDependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res := installer.FixSingleDependency(req)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleFixIssue(w http.ResponseWriter, r *http.Request) {
	var req installer.FixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res := installer.FixIssue(req)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handlePackageUninstall(w http.ResponseWriter, r *http.Request) {
	var req installer.UninstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res := installer.UninstallPackage(req)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handlePackageSwitchVersion(w http.ResponseWriter, r *http.Request) {
	var req installer.SwitchVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res := installer.SwitchPackageVersion(req)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleProfileExecuteBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Commands []string `json:"commands"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res := installer.ExecuteBatchCommands(req.Commands)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleProfileExportZip(w http.ResponseWriter, r *http.Request) {
	var req profile.SelectiveExportRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	state, err := scanner.ScanAll(s.scannerOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zipBytes, err := profile.ExportZipBundle(state, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=packetinstall-bundle.zip")
	_, _ = w.Write(zipBytes)
}

func (s *Server) handleProfileExportSelective(w http.ResponseWriter, r *http.Request) {
	var req profile.SelectiveExportRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	state, err := scanner.ScanAll(s.scannerOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	yamlBytes, err := profile.ExportCustomYAML(state, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", "attachment; filename=packetinstall.yaml")
	_, _ = w.Write(yamlBytes)
}

func (s *Server) handleProjectsScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path  string `json:"path"`
		Depth int    `json:"depth"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		req.Path = "D:/IdeaSideProject"
	}
	if req.Depth <= 0 {
		req.Depth = 3
	}

	result, err := scanner.ScanProjects(req.Path, req.Depth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (s *Server) handlePackageInstall(w http.ResponseWriter, r *http.Request) {
	var req model.InstallPackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res := installer.InstallPackage(req)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleAdapters(w http.ResponseWriter, r *http.Request) {
	type Adapter struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"` // Available, Fixture only, Fixture unavailable
	}

	adapters := []Adapter{
		{Name: "Chocolatey", Description: "Windows package manager & tool inventory", Status: "Available"},
		{Name: "npm", Description: "User-scoped global CLI ownership", Status: "Available"},
		{Name: "Scoop", Description: "User-space Windows command line installers", Status: "Available"},
		{Name: "MCP client configs", Description: "Codex, Claude Code, and Cursor global servers", Status: "Available"},
		{Name: "WinGet", Description: "Windows Package Manager mappings", Status: "Available"},
		{Name: "Linux native manager", Description: "APT, DNF, or Pacman", Status: "Fixture unavailable"},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(adapters)
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	type HistoryItem struct {
		Status    string `json:"status"` // Success, Partial, Cancelled, Recoverable, Failed
		Resource  string `json:"resource"`
		Action    string `json:"action"`
		Timestamp string `json:"timestamp"`
	}

	history := []HistoryItem{
		{Status: "Success", Resource: "Codex CLI", Action: "Update", Timestamp: "Just now"},
		{Status: "Success", Resource: "Git", Action: "Inventory refresh", Timestamp: "2 mins ago"},
		{Status: "Partial", Resource: "Browser Control", Action: "Update skill", Timestamp: "10 mins ago"},
		{Status: "Recoverable", Resource: "Release Pilot", Action: "Update skill", Timestamp: "25 mins ago"},
		{Status: "Cancelled", Resource: "Docker Desktop", Action: "Vendor handoff", Timestamp: "1 hour ago"},
		{Status: "Success", Resource: "fable-thinking", Action: "Catalog import", Timestamp: "2 hours ago"},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(history)
}

func (s *Server) handleUpdates(w http.ResponseWriter, r *http.Request) {
	type UpdateQueueItem struct {
		ID        string `json:"id"`
		Type      string `json:"type"` // tool, skill
		Name      string `json:"name"`
		Version   string `json:"version"` // e.g. "0.9.4 -> 0.10.1"
		Authority string `json:"authority"` // vendor handoff, managed execute
		Risk      string `json:"risk"` // Vendor execution, 2 target files changed, etc.
	}

	updates := []UpdateQueueItem{
		{ID: "claude-code", Type: "tool", Name: "Claude Code", Version: "2.1.246 -> 2.1.261", Authority: "managed execute", Risk: "Upstream patch update"},
		{ID: "docker", Type: "tool", Name: "Docker Desktop", Version: "4.44.2 -> 4.45.0", Authority: "vendor handoff", Risk: "Vendor execution"},
		{ID: "frontend-design", Type: "skill", Name: "Frontend Design", Version: "7f84c21 -> c9e9f31", Authority: "managed execute", Risk: "2 target files changed"},
		{ID: "release-pilot", Type: "skill", Name: "Release Pilot", Version: "d24b80c -> f91a6bc", Authority: "managed execute", Risk: "Blocked by local modification"},
		{ID: "browser-control", Type: "skill", Name: "Browser Control", Version: "v0.8.0 -> v0.9.0", Authority: "managed execute", Risk: "Dependency conflict"},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updates)
}

func (s *Server) handlePackageDiagnose(w http.ResponseWriter, r *http.Request) {
	var pkg model.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	diag := s.auditor.DiagnosePackage(pkg)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(diag)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	state, err := scanner.ScanAll(s.scannerOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(state)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	state, err := scanner.ScanAll(s.scannerOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []model.AuditItem
	for _, pkg := range state.Packages {
		if pkg.Manager == "npm" {
			if item, err := s.auditor.CheckNpmPackage(pkg.Name, pkg.Version); err == nil {
				items = append(items, *item)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func (s *Server) handleProfileExport(w http.ResponseWriter, r *http.Request) {
	state, err := scanner.ScanAll(s.scannerOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	yamlBytes, err := profile.ExportProfileYAML(state, "my-machine-profile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", "attachment; filename=packetinstall.yaml")
	_, _ = w.Write(yamlBytes)
}

func (s *Server) handleProfileDiff(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	targetProfile, err := profile.ImportProfileYAML(body)
	if err != nil {
		http.Error(w, "invalid YAML profile: "+err.Error(), http.StatusBadRequest)
		return
	}

	currentState, err := scanner.ScanAll(s.scannerOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	diff := profile.CalculateDiff(currentState, targetProfile)
	cmds := profile.GenerateInstallPlan(diff, runtime.GOOS, "")

	resp := map[string]interface{}{
		"already_installed":        diff.AlreadyInstalled,
		"pending_system_packages": diff.PendingSystemPackages,
		"pending_global_clis":     diff.PendingGlobalCLIs,
		"missing_skills":          diff.MissingSkills,
		"missing_mcp_servers":     diff.MissingMcpServers,
		"commands":                cmds,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	data, err := staticFS.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "static dashboard not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}
