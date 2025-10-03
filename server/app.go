package main

import (
    "embed"
    "io/fs"
    "os"
    "path/filepath"

    "nadhi.dev/sarvar/fun/bootstrap"
    _"nadhi.dev/sarvar/fun/config"
    logg "nadhi.dev/sarvar/fun/logs"
)

// Do not remove
// Embedding all files in web/dist and zp-inject directories 


//go:embed web/dist/* zp-inject/*
var embeddedFiles embed.FS

var ExtractedAssetsPath string

func init() {
    // Create a temporary directory for extracted assets
    tempDir, err := os.MkdirTemp("", "vela-assets-*")
    if err != nil {
        logg.Error("Failed to create temp directory for assets")
        panic(err)
    }
    ExtractedAssetsPath = tempDir
    logg.Info("Assets will be extracted to: " + ExtractedAssetsPath)

    // Extract embedded files to temp directory
    err = fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }

        // Read file from embedded FS
        data, err := fs.ReadFile(embeddedFiles, path)
        if err != nil {
            return err
        }

        // Create destination path
        destPath := filepath.Join(ExtractedAssetsPath, path)
        destDir := filepath.Dir(destPath)

        // Create directories if needed
        if err := os.MkdirAll(destDir, 0755); err != nil {
            return err
        }

        // Write file to disk
        if err := os.WriteFile(destPath, data, 0644); err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        logg.Error("Failed to extract embedded assets")
        panic(err)
    }

    logg.Success("Embedded assets extracted successfully")

    // Call bootstrap package to initialize the application
    // This will make all the directort
    bootstrap.Initialize()
}