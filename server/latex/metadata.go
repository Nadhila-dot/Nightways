package latex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetMetadataForJob(texFilename string) (map[string]interface{}, error) {
	metaFilename := strings.TrimSuffix(texFilename, filepath.Ext(texFilename)) + ".meta.json"
	sourceFile := metaFilename

	if _, err := os.Stat(sourceFile); err != nil {
		possiblePaths := []string{
			filepath.Join(".", metaFilename),
			filepath.Join("./generated", filepath.Base(metaFilename)),
		}

		parts := strings.Split(filepath.Base(texFilename), ".")
		if len(parts) > 0 {
			jobID := parts[0]
			possiblePaths = append(possiblePaths,
				filepath.Join("./generated", jobID, filepath.Base(metaFilename)))
		}

		found := false
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				sourceFile = path
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("metadata file not found for %s", texFilename)
		}
	}

	metadataBytes, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %v", err)
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %v", err)
	}

	return metadata, nil
}

func UploadMetadataToStorage(texFilename string, bucketURL string) (map[string]interface{}, string, error) {
	metadata, err := GetMetadataForJob(texFilename)
	if err != nil {
		return nil, "", err
	}

	metadataFilename := strings.TrimSuffix(filepath.Base(texFilename), filepath.Ext(texFilename)) + ".meta.json"
	metadataURL := fmt.Sprintf("%s/%s", bucketURL, metadataFilename)

	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal metadata: %v", err)
	}

	tempDir, err := ioutil.TempDir("", "metadata")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempPath := filepath.Join(tempDir, metadataFilename)
	if err := ioutil.WriteFile(tempPath, metadataBytes, 0644); err != nil {
		return nil, "", fmt.Errorf("failed to write metadata to temp file: %v", err)
	}

	// Example: later we add actual upload logic here
	// err = UploadToStorage(tempPath, metadataURL)

	return metadata, metadataURL, nil
}
