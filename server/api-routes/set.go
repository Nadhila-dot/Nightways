package api

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
   _ "os"
    "strings"
)

// Auth check: expects "Authorization: Bearer <session>"
func isAuthorized(r *http.Request) bool {
    auth := r.Header.Get("Authorization")
    return strings.HasPrefix(auth, "Bearer ") && len(auth) > 7
}

// POST /api/set
// Body: JSON object to replace set.json
func SetSystemEnv(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if !isAuthorized(r) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var newData interface{}
    if err := json.NewDecoder(r.Body).Decode(&newData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Write to set.json
    filePath := "./server/set.json"
    dataBytes, err := json.MarshalIndent(newData, "", "  ")
    if err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        return
    }
    if err := ioutil.WriteFile(filePath, dataBytes, 0644); err != nil {
        http.Error(w, "Failed to write file", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message":"System environment updated successfully"}`))
}

func GetSystemEnv(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    filePath := "./server/set.json"
    dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
        http.Error(w, "Failed to read file", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(dataBytes)
}