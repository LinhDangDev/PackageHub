package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"packetinstall/internal/model"
)

type parsedSkillMeta struct {
	Name         string
	Description  string
	Category     string
	ArgumentHint string
	WhenToUse    string
	Author       string
	FullDocs     string
}

// ScanSkills scans directories containing AI agent skills and reads their SKILL.md.
func ScanSkills(baseDir string) ([]model.Skill, error) {
	var skills []model.Skill
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return skills, nil
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || entry.Name() == "_shared" {
			continue
		}
		skillDir := filepath.Join(baseDir, entry.Name())
		skillMdPath := filepath.Join(skillDir, "SKILL.md")

		meta := parsedSkillMeta{
			Name: entry.Name(),
		}

		if data, err := os.ReadFile(skillMdPath); err == nil {
			meta = parseSkillMd(data, entry.Name())
		}

		gitRemote, commitSHA := getGitRepoInfo(skillDir)
		toolOrigin := determineToolOrigin(meta.Name, meta.Author, skillDir)
		tier := determineSkillTier(meta.Name)
		examples := generateSkillExamples(meta.Name, meta.ArgumentHint)
		flags := extractSkillFlags(meta.ArgumentHint, meta.FullDocs)

		cmd := meta.Name
		if !strings.HasPrefix(cmd, "/") {
			cmd = "/" + cmd
		}

		skills = append(skills, model.Skill{
			Name:         meta.Name,
			ToolOrigin:   toolOrigin,
			Category:     meta.Category,
			Tier:         tier,
			Command:      cmd,
			ArgumentHint: meta.ArgumentHint,
			WhenToUse:    meta.WhenToUse,
			Path:         skillDir,
			Description:  meta.Description,
			GitRemote:    gitRemote,
			CommitSHA:    commitSHA,
			FullDocs:     meta.FullDocs,
			Examples:     examples,
			Flags:        flags,
		})
	}

	return skills, nil
}

func parseSkillMd(content []byte, folderName string) parsedSkillMeta {
	meta := parsedSkillMeta{
		Name: folderName,
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	inFrontmatter := false
	var bodyBuilder strings.Builder
	pastFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "---" {
			if !inFrontmatter && !pastFrontmatter {
				inFrontmatter = true
				continue
			} else if inFrontmatter {
				inFrontmatter = false
				pastFrontmatter = true
				continue
			}
		}

		if inFrontmatter {
			if strings.HasPrefix(trimmed, "name:") {
				meta.Name = cleanYamlValue(strings.TrimPrefix(trimmed, "name:"))
			} else if strings.HasPrefix(trimmed, "description:") {
				meta.Description = cleanYamlValue(strings.TrimPrefix(trimmed, "description:"))
			} else if strings.HasPrefix(trimmed, "category:") {
				meta.Category = cleanYamlValue(strings.TrimPrefix(trimmed, "category:"))
			} else if strings.HasPrefix(trimmed, "argument-hint:") {
				meta.ArgumentHint = cleanYamlValue(strings.TrimPrefix(trimmed, "argument-hint:"))
			} else if strings.HasPrefix(trimmed, "when_to_use:") {
				meta.WhenToUse = cleanYamlValue(strings.TrimPrefix(trimmed, "when_to_use:"))
			} else if strings.HasPrefix(trimmed, "author:") {
				meta.Author = cleanYamlValue(strings.TrimPrefix(trimmed, "author:"))
			}
		} else if pastFrontmatter {
			bodyBuilder.WriteString(line + "\n")
		}
	}

	meta.FullDocs = strings.TrimSpace(bodyBuilder.String())
	if meta.Category == "" {
		meta.Category = "utilities"
	}
	return meta
}

func cleanYamlValue(val string) string {
	val = strings.TrimSpace(val)
	val = strings.Trim(val, "\"'")
	return val
}

func determineToolOrigin(name, author, dir string) string {
	lowerName := strings.ToLower(name)
	lowerAuthor := strings.ToLower(author)
	lowerDir := strings.ToLower(dir)

	if strings.HasPrefix(lowerName, "ck:") || lowerAuthor == "claudekit" || strings.Contains(lowerDir, "claudekit") {
		return "claudekit"
	}
	if strings.Contains(lowerDir, ".omp") || strings.HasPrefix(lowerName, "omp:") {
		return "opencode"
	}
	if strings.Contains(lowerDir, ".codex") || strings.HasPrefix(lowerName, "codex:") || strings.HasPrefix(lowerName, "codex-") {
		return "codex"
	}
	if strings.Contains(lowerDir, "cursor") {
		return "cursor"
	}
	return "universal"
}

func determineSkillTier(name string) string {
	n := strings.TrimPrefix(name, "ck:")
	switch n {
	case "ask", "brainstorm", "fix", "docs", "interview-docs", "find-skills", "coding-level":
		return "beginner"
	case "plan", "cook", "test", "review-pr", "git", "preview", "ui-styling", "frontend-development", "backend-development", "databases":
		return "intermediate"
	case "bootstrap", "scout", "ship", "watzup", "team", "vibe", "autoresearch", "loop", "security", "cti-expert", "mcp-builder", "ui-ux-pro-max", "ai-artist":
		return "pro"
	default:
		return "intermediate"
	}
}

func generateSkillExamples(name, argHint string) []string {
	n := strings.TrimPrefix(name, "ck:")
	switch n {
	case "ask":
		return []string{
			`/ck:ask "How does authentication work in this app?"`,
			`/ck:ask "Evaluate trade-offs between PostgreSQL and SQLite"`,
		}
	case "plan":
		return []string{
			`/ck:plan "Add user profile page" --tdd`,
			`/ck:plan "Migrate auth service" --fast`,
		}
	case "cook":
		return []string{
			`/ck:cook "Add dark mode toggle to settings" --fast`,
			`/ck:cook @plan.md --parallel`,
		}
	case "brainstorm":
		return []string{
			`/ck:brainstorm "Real-time notification architecture"`,
			`/ck:brainstorm "Mobile UX onboarding flow"`,
		}
	case "fix":
		return []string{
			`/ck:fix "500 error in /api/users endpoint"`,
			`/ck:fix "Fix failing unit tests in payment service"`,
		}
	case "scout":
		return []string{
			`/ck:scout "Locate all payment webhook handlers"`,
			`/ck:scout "Find database migration files"`,
		}
	case "ship":
		return []string{
			`/ck:ship "Prepare v1.2.0 production release"`,
			`/ck:ship --auto`,
		}
	case "watzup":
		return []string{
			`/ck:watzup`,
			`/ck:watzup --diff`,
		}
	case "ui-ux-pro-max":
		return []string{
			`/ck:ui-ux-pro-max "Design modern dashboard layout with Tailwind"`,
			`/ck:ui-ux-pro-max "Audit accessibility and micro-interactions"`,
		}
	case "ai-artist":
		return []string{
			`/ck:ai-artist "product showcase in cyberpunk style" --mode creative`,
			`/ck:ai-artist "tech conference banner" --provider google --skip`,
		}
	default:
		cmd := "/" + name
		if argHint != "" {
			return []string{
				fmt.Sprintf("%s %s", cmd, argHint),
			}
		}
		return []string{cmd}
	}
}

var standardFlagDesc = map[string]string{
	"--fast":        "Bỏ qua bước nghiên cứu sơ bộ, chuyển thẳng từ scout ➔ plan ➔ sinh code ngay.",
	"--parallel":    "Kích hoạt điều phối đa agent chạy song song để tăng tốc độ thực thi.",
	"--tdd":         "Quy trình Test-Driven Development: viết test trước khi viết code, bảo đảm pass 100%.",
	"--auto":        "Tự động phê duyệt toàn bộ các bước mà không cần dừng lại hỏi xác nhận.",
	"--no-test":     "Bỏ qua bước chạy kiểm thử (dùng khi cần sửa nhanh tài liệu hoặc UI).",
	"--advice":      "Bật giám sát cố vấn từ agent kongming (mô hình suy luận cao cấp) để tư vấn kiến trúc.",
	"--interactive": "Chế độ tương tác đầy đủ: dừng lại hỏi ý kiến người dùng tại từng bước.",
	"--hard":        "Chế độ thực thi nghiêm ngặt: yêu cầu kiểm tra toàn diện và review phản biện.",
	"--deep":        "Phân tích và lập kế hoạch chuyên sâu, đào sâu vào kiến trúc và biên độ rủi ro.",
	"--two":         "Tạo 2 kế hoạch giải pháp cạnh tranh nhau để người dùng lựa chọn hướng đi tốt nhất.",
	"--html":        "Xuất bản tài liệu/kế hoạch thành file HTML độc lập có giao diện trực quan tuyệt đẹp.",
	"--github":      "Tự động tạo Issue hoặc PR trên GitHub và gán nhãn ready to review.",
	"--wiki":        "Xuất bản tài liệu trực tiếp lên AgentWiki.",
	"--skip":        "Bỏ qua bước phỏng vấn bắt buộc ban đầu (bỏ qua interview style/mood/colors).",
	"--mode":        "Lựa chọn chế độ sinh ảnh: search (mẫu có sẵn), creative (kết hợp), wild (sáng tạo tự do), all (cả 3).",
	"--provider":    "Lựa chọn nhà cung cấp AI: auto (tự động), google (Gemini API), openrouter.",
	"--explore":     "Mở rộng phạm vi tìm kiếm codebase sang các thư mục liên quan.",
	"--native":      "Sử dụng công cụ tìm kiếm nội tại thay vì công cụ agent mở rộng.",
	"--fix":         "Tự động sửa các lỗi hoặc nhận xét tìm thấy trong quá trình review PR.",
	"--merge":       "Tự động merge Pull Request sau khi CI và review vượt qua.",
	"--reply":       "Gửi trực tiếp nhận xét review lên GitHub PR qua gh CLI.",
	"--diff":        "Hiển thị chi tiết khác biệt git giữa các commit hoặc nhánh.",
	"--branch":      "Chỉ định nhánh git cụ thể cần thao tác.",
	"--verbose":     "Hiển thị thông tin chi tiết đầy đủ trong quá trình chạy.",
	"--no-tasks":    "Bỏ qua việc đồng bộ task lên hệ thống task manager runtime.",
}

func extractSkillFlags(argHint, fullDocs string) []model.SkillFlag {
	flagsMap := make(map[string]model.SkillFlag)
	var orderedNames []string

	// 1. Parse from markdown fullDocs bullet points: e.g. "- `--fast`: Skip research..."
	scanner := bufio.NewScanner(strings.NewReader(fullDocs))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "- `--") || strings.HasPrefix(line, "* `--") || strings.HasPrefix(line, "- --") {
			clean := strings.TrimPrefix(line, "- ")
			clean = strings.TrimPrefix(clean, "* ")
			parts := strings.SplitN(clean, ":", 2)
			if len(parts) == 2 {
				flagName := strings.Trim(parts[0], "`' \"")
				desc := strings.TrimSpace(parts[1])
				desc = strings.Trim(desc, "`'\"")
				if strings.HasPrefix(flagName, "--") {
					flagsMap[flagName] = model.SkillFlag{
						Name:        flagName,
						Description: desc,
					}
					orderedNames = append(orderedNames, flagName)
				}
			}
		}
	}

	// 2. Parse from argument-hint tokens: e.g. "[concept] [--mode search|creative|wild|all] [--skip]"
	rawTokens := strings.Fields(argHint)
	for _, tok := range rawTokens {
		tok = strings.Trim(tok, "[]()")
		if strings.Contains(tok, "--") {
			// Handle pipe-separated flags, e.g. "--fast|--hard|--deep"
			subTokens := strings.Split(tok, "|")
			for _, sub := range subTokens {
				sub = strings.Trim(sub, "[]()")
				if strings.HasPrefix(sub, "--") {
					parts := strings.Fields(sub)
					fName := parts[0]
					val := ""
					if len(parts) > 1 {
						val = strings.Join(parts[1:], " ")
					}

					if _, exists := flagsMap[fName]; !exists {
						desc := standardFlagDesc[fName]
						if desc == "" {
							desc = fmt.Sprintf("Tùy chọn cờ điều khiển %s cho lệnh.", fName)
						}
						flagsMap[fName] = model.SkillFlag{
							Name:        fName,
							Values:      val,
							Description: desc,
						}
						orderedNames = append(orderedNames, fName)
					} else if val != "" {
						f := flagsMap[fName]
						f.Values = val
						flagsMap[fName] = f
					}
				}
			}
		}
	}

	// Build ordered slice
	result := make([]model.SkillFlag, 0, len(orderedNames))
	seen := make(map[string]bool)
	for _, name := range orderedNames {
		if !seen[name] {
			seen[name] = true
			if f, ok := flagsMap[name]; ok {
				if std, ok := standardFlagDesc[name]; ok && (len(f.Description) < 15 || !strings.Contains(f.Description, " ")) {
					f.Description = std
				}
				result = append(result, f)
			}
		}
	}

	return result
}

// getGitRepoInfo parses .git/config and .git/HEAD directly without running git.exe or spawning console windows.
func getGitRepoInfo(dir string) (string, string) {
	gitDir := filepath.Join(dir, ".git")
	fi, err := os.Stat(gitDir)
	if err != nil {
		return "", ""
	}

	// Handle git submodules or worktrees where .git is a file containing "gitdir: ..."
	if !fi.IsDir() {
		data, err := os.ReadFile(gitDir)
		if err == nil {
			line := strings.TrimSpace(string(data))
			if strings.HasPrefix(line, "gitdir:") {
				relPath := strings.TrimSpace(strings.TrimPrefix(line, "gitdir:"))
				gitDir = filepath.Join(dir, relPath)
			}
		}
	}

	remote := ""
	cfgData, err := os.ReadFile(filepath.Join(gitDir, "config"))
	if err == nil {
		lines := strings.Split(string(cfgData), "\n")
		inOrigin := false
		for _, l := range lines {
			trimmed := strings.TrimSpace(l)
			if strings.HasPrefix(trimmed, "[remote \"origin\"]") {
				inOrigin = true
				continue
			} else if inOrigin && strings.HasPrefix(trimmed, "[") {
				inOrigin = false
			}
			if inOrigin && strings.HasPrefix(trimmed, "url =") {
				remote = strings.TrimSpace(strings.TrimPrefix(trimmed, "url ="))
				break
			}
		}
	}

	commit := ""
	headData, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err == nil {
		headContent := strings.TrimSpace(string(headData))
		if strings.HasPrefix(headContent, "ref:") {
			refPath := strings.TrimSpace(strings.TrimPrefix(headContent, "ref:"))
			refFile := filepath.Join(gitDir, refPath)
			if refData, err := os.ReadFile(refFile); err == nil {
				hash := strings.TrimSpace(string(refData))
				if len(hash) >= 7 {
					commit = hash[:7]
				}
			} else {
				if packedData, err := os.ReadFile(filepath.Join(gitDir, "packed-refs")); err == nil {
					pLines := strings.Split(string(packedData), "\n")
					for _, pl := range pLines {
						if strings.HasSuffix(pl, refPath) {
							parts := strings.Fields(pl)
							if len(parts) > 0 && len(parts[0]) >= 7 {
								commit = parts[0][:7]
								break
							}
						}
					}
				}
			}
		} else if len(headContent) >= 7 {
			commit = headContent[:7]
		}
	}

	return remote, commit
}
