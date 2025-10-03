package bootstrap

import (
    logg "nadhi.dev/sarvar/fun/logs"
)

// Initialize sets up all required components for the application
func Initialize() {
    logg.Info("Initializing application...")
    
    // Initialize configs first
    InitConfigs()
    
    // Initialize database structure
   // InitDatabase()
    
    logg.Success("Application initialization complete")
}