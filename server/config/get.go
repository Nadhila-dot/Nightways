package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

func getConfigPath() string {
    // Always use the same logic as SaveConfig
    rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(os.Args[0]))))
    return filepath.Join(rootDir, "set.json")
}

func GetConfigValue(key string) interface{} {
    configPath := getConfigPath()
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil
    }
    var arr []map[string]interface{}
    if err := json.Unmarshal(data, &arr); err != nil || len(arr) == 0 {
        return nil
    }
    return arr[0][key]
}

func GetConfig() map[string]interface{} {
    configPath := getConfigPath()
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil
    }
    var config map[string]interface{}
    if err := json.Unmarshal(data, &config); err != nil {
        return nil
    }
    return config
}