package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func MakeJsonConfig(api_key string, source string, model string) error {
    config := map[string]interface{}{
        "AI_API": api_key,
        "MAX_SESSIONS":     2,
        "SHEET_QUEUE_DIR":  "./storage/queue_data",
    }
    return SaveConfig(config)
}

func SaveConfig(config map[string]interface{}) error {
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(os.Args[0]))))
	configPath := filepath.Join(rootDir, "set.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
