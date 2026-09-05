package scanner

import (
	"encoding/json"
	"os"

	"packetinstall/internal/model"
)

type mcpConfigFile struct {
	McpServers map[string]struct {
		Command  string            `json:"command"`
		Args     []string          `json:"args"`
		Env      map[string]string `json:"env"`
		Disabled bool              `json:"disabled"`
	} `json:"mcpServers"`
}

// ScanMcpConfigFile parses Claude Desktop or Cursor mcpServers JSON configuration.
func ScanMcpConfigFile(filePath string, source string) ([]model.McpServer, error) {
	var servers []model.McpServer
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return servers, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg mcpConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	for name, s := range cfg.McpServers {
		servers = append(servers, model.McpServer{
			Name:     name,
			Source:   source,
			Command:  s.Command,
			Args:     s.Args,
			Env:      s.Env,
			Disabled: s.Disabled,
		})
	}

	return servers, nil
}
