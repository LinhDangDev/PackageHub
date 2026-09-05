package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"packetinstall/internal/app"
	"packetinstall/internal/auditor"
	"packetinstall/internal/profile"
	"packetinstall/internal/scanner"
	"packetinstall/internal/web"
)

const version = "1.2.0"

func attachParentConsole() {
	if runtime.GOOS == "windows" {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		attachConsole := kernel32.NewProc("AttachConsole")
		// ATTACH_PARENT_PROCESS is (DWORD)-1 = 0xFFFFFFFF
		r, _, _ := attachConsole.Call(uintptr(0xFFFFFFFF))
		if r != 0 {
			os.Stdout = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
			os.Stderr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
		}
	}
}

func printBanner() {
	fmt.Println("⚡ packetinstall v" + version + " — Smart Tools, Skills & Agent Profile Manager")
}

func printUsage() {
	printBanner()
	fmt.Println(`
Usage: packetinstall [command] [options]

Commands:
  app                 Launch native Windows Desktop Application window (Default)
  scan                Fast-scan all packages (Choco, Scoop, NPM), Agent Skills & MCP servers
  audit               Audit installed tools for End-Of-Life (EOL) runtimes and updates
  export [-o file]    Export current environment profile into a declarative YAML manifest
  apply <file>        Calculate diff and apply/install profile on this machine (--dry-run supported)
  ui [-port 3456]     Launch embedded Web Dashboard in your browser
  version             Display version information

Examples:
  packetinstall                                  # Opens native Windows Desktop App (zero console window)
  packetinstall scan                             # Fast CLI terminal scan
  packetinstall audit                            # Check for outdated / EOL runtimes
  packetinstall export -o my-profile.yaml         # Export current machine profile
  packetinstall apply my-profile.yaml --dry-run   # Preview diff and install commands`)
}

func main() {
	// If run without arguments, or with "app", launch the Native Windows Desktop App!
	if len(os.Args) < 2 || os.Args[1] == "app" {
		opts := scanner.DefaultScanOptions()
		if err := app.RunDesktopApp(opts); err != nil {
			attachParentConsole()
			fmt.Printf("⚠️ Could not open native window: %v\nFalling back to browser UI...\n", err)
			runUI([]string{})
		}
		return
	}

	// Attach to caller console for CLI commands
	attachParentConsole()

	command := os.Args[1]

	switch command {
	case "scan":
		runScan()
	case "audit":
		runAudit()
	case "export":
		runExport(os.Args[2:])
	case "apply":
		runApply(os.Args[2:])
	case "ui":
		runUI(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	case "version", "--version", "-v":
		fmt.Println("packetinstall v" + version)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runScan() {
	printBanner()
	fmt.Println("🔍 Scanning system (Choco, Scoop, NPM Globals, Skills, MCP configs)...")

	opts := scanner.DefaultScanOptions()
	state, err := scanner.ScanAll(opts)
	if err != nil {
		fmt.Printf("❌ Scan failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n📦 Installed Packages & CLIs (%d detected):\n", len(state.Packages))
	fmt.Printf("  %-8s %-35s %-15s\n", "MANAGER", "NAME", "VERSION")
	fmt.Println("  " + strings.Repeat("-", 60))
	for _, p := range state.Packages {
		v := p.Version
		if v == "" {
			v = "unknown"
		}
		fmt.Printf("  %-8s %-35s %-15s\n", p.Manager, p.Name, v)
	}

	fmt.Printf("\n🧠 AI Agent Skills (%d detected):\n", len(state.Skills))
	for _, s := range state.Skills {
		desc := s.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("  • %-25s — %s\n", s.Name, desc)
	}

	fmt.Printf("\n🔌 MCP Servers (%d configured):\n", len(state.McpServers))
	for _, m := range state.McpServers {
		fmt.Printf("  • [%s] %-15s -> %s %s\n", m.Source, m.Name, m.Command, strings.Join(m.Args, " "))
	}

	fmt.Printf("\n⚡ Scan complete in %dms!\n", state.DurationMs)
}

func runAudit() {
	printBanner()
	fmt.Println("🔍 Auditing tools against End-Of-Life (EOL) database and upstream registries...")

	opts := scanner.DefaultScanOptions()
	state, err := scanner.ScanAll(opts)
	if err != nil {
		fmt.Printf("❌ Scan failed: %v\n", err)
		os.Exit(1)
	}

	client := auditor.NewAuditorClient()

	fmt.Println("\n📊 Health & Obsolescence Report:")
	fmt.Printf("  %-30s %-12s %-12s %-15s %s\n", "PACKAGE / RUNTIME", "CURRENT", "LATEST", "STATUS", "ACTION")
	fmt.Println("  " + strings.Repeat("-", 85))

	for _, p := range state.Packages {
		if p.Manager == "npm" {
			item, err := client.CheckNpmPackage(p.Name, p.Version)
			if err == nil && item != nil {
				statusColor := item.Status
				fmt.Printf("  %-30s %-12s %-12s %-15s %s\n", item.Name, item.CurrentVersion, item.LatestVersion, statusColor, item.UpdateCommand)
			}
		}
	}
	fmt.Println("\n✅ Audit complete!")
}

func runExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	outFile := fs.String("o", "packetinstall.yaml", "Output YAML profile path")
	_ = fs.Parse(args)

	printBanner()
	fmt.Println("📦 Exporting current machine profile...")

	opts := scanner.DefaultScanOptions()
	state, err := scanner.ScanAll(opts)
	if err != nil {
		fmt.Printf("❌ Failed to scan system: %v\n", err)
		os.Exit(1)
	}

	yamlBytes, err := profile.ExportProfileYAML(state, "exported-profile")
	if err != nil {
		fmt.Printf("❌ Failed to generate profile: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outFile, yamlBytes, 0644); err != nil {
		fmt.Printf("❌ Failed to write file %s: %v\n", *outFile, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Profile exported successfully to: %s\n", *outFile)
	fmt.Printf("   • Packages: %d\n   • Skills: %d\n   • MCP Servers: %d (Secrets safely masked)\n", len(state.Packages), len(state.Skills), len(state.McpServers))
}

func runApply(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: packetinstall apply <profile.yaml> [--dry-run]")
		os.Exit(1)
	}

	profileFile := args[0]
	dryRun := false
	for _, a := range args[1:] {
		if a == "--dry-run" {
			dryRun = true
		}
	}

	printBanner()
	fmt.Printf("📥 Reading profile: %s\n", profileFile)

	data, err := os.ReadFile(profileFile)
	if err != nil {
		fmt.Printf("❌ Cannot read file: %v\n", err)
		os.Exit(1)
	}

	targetProfile, err := profile.ImportProfileYAML(data)
	if err != nil {
		fmt.Printf("❌ Invalid profile format: %v\n", err)
		os.Exit(1)
	}

	opts := scanner.DefaultScanOptions()
	currentState, err := scanner.ScanAll(opts)
	if err != nil {
		fmt.Printf("❌ Cannot scan current machine: %v\n", err)
		os.Exit(1)
	}

	diff := profile.CalculateDiff(currentState, targetProfile)
	commands := profile.GenerateInstallPlan(diff, runtime.GOOS, opts.SkillsDirs[0])

	fmt.Printf("\n📋 Diff Analysis vs Current Machine:\n")
	fmt.Printf("  • Already Installed:      %d items (skipped)\n", len(diff.AlreadyInstalled))
	fmt.Printf("  • Missing System Packages: %d\n", len(diff.PendingSystemPackages))
	for mgr, pkgs := range diff.PendingGlobalCLIs {
		fmt.Printf("  • Missing %s Global CLIs: %d\n", mgr, len(pkgs))
	}
	fmt.Printf("  • Missing Skills:          %d\n", len(diff.MissingSkills))
	fmt.Printf("  • Missing MCP Servers:     %d\n", len(diff.MissingMcpServers))

	if len(commands) == 0 {
		fmt.Println("\n🎉 Machine is already fully synced with this profile!")
		return
	}

	fmt.Printf("\n🚀 Generated Silent Installation Plan (%d commands):\n", len(commands))
	for i, cmd := range commands {
		fmt.Printf("  [%d] %s\n", i+1, cmd)
	}

	if dryRun {
		fmt.Println("\n🔍 [Dry-Run Mode] No changes were made to your system.")
		return
	}

	fmt.Println("\n⚡ To execute installation, run each command or pipe into your terminal.")
}

func runUI(args []string) {
	fs := flag.NewFlagSet("ui", flag.ExitOnError)
	port := fs.Int("port", 3456, "Local server port")
	_ = fs.Parse(args)

	printBanner()
	opts := scanner.DefaultScanOptions()
	server := web.NewServer(opts)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	url := fmt.Sprintf("http://localhost:%d", *port)

	fmt.Printf("🚀 Starting Web Dashboard at %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")

	go func() {
		if runtime.GOOS == "windows" {
			_ = exec.Command("cmd", "/c", "start", url).Start()
		} else if runtime.GOOS == "darwin" {
			_ = exec.Command("open", url).Start()
		}
	}()

	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		fmt.Printf("❌ Server error: %v\n", err)
	}
}
