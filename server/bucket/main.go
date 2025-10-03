package vela

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

// List all files in ./storage and send as JSON array
func ListStorageFiles(w http.ResponseWriter, r *http.Request) {
    files := []string{}
    root := "./storage"

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            relPath := strings.TrimPrefix(path, root+"/")
            files = append(files, relPath)
        }
        return nil
    })
    if err != nil {
        http.Error(w, "Failed to list files", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}

// Serve a file from ./storage/bucket/whatever
func ServeStorageFile(w http.ResponseWriter, r *http.Request) {
    // Example: /vela/bucket/image.png
    // Get everything after /vela/bucket/
    relPath := strings.TrimPrefix(r.URL.Path, "/vela/bucket/")
    filePath := filepath.Join("./storage", relPath)

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Set content type based on file extension (simple)
    switch ext := strings.ToLower(filepath.Ext(filePath)); ext {
    case ".png":
        w.Header().Set("Content-Type", "image/png")
    case ".jpg", ".jpeg":
        w.Header().Set("Content-Type", "image/jpeg")
    case ".pdf":
        w.Header().Set("Content-Type", "application/pdf")
    default:
        w.Header().Set("Content-Type", "application/octet-stream")
    }

    w.Write(data)
}