package bootstrap

import (
   _ "encoding/json"
    "os"
   _ "path/filepath"

    "nadhi.dev/sarvar/fun/config"
    logg "nadhi.dev/sarvar/fun/logs"
)

func getConfigPath() string {
    return "./set.json"  // Look in the current working directory
}

// InitConfigs ensures that the set.json configuration file exists
// and has the required structure
func InitConfigs() {
    configPath := "./set.json"
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        logg.Warning("set.json not found, creating default config...")
        
        // Create a default config file
        defaultConfig := map[string]interface{}{
            "AI_API":          "",
            "MAX_SESSIONS":    2,
            "SHEET_QUEUE_DIR": "./storage/queue_data",
        }
        
        if err := config.SaveConfig(defaultConfig); err != nil {
            logg.Error("Failed to create default config: " + err.Error())
            logg.Exit()
        }
        
        logg.Success("Default set.json created at " + configPath)
    }
}