package kasm

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
	"github.com/joho/godotenv"
)

type TargetUser struct {
    UserID string `json:"user_id"`
}

type LoginRequest struct {
    APIKey       string     `json:"api_key"`
    APIKeySecret string     `json:"api_key_secret"`
    TargetUser   TargetUser `json:"target_user"`
}

type LoginResponse struct {
    URL string `json:"url"`
}



func GetLoginLink(userID string) (string, error) {
	// Load API credentials from .env file
	_ = godotenv.Load(".env")
	apiKey := os.Getenv("API_KEY")
	apiKeySecret := os.Getenv("API_KEY_SECRET")
	if apiKey == "" || apiKeySecret == "" {
		return "", fmt.Errorf("API credentials not set")
	}

	reqBody := LoginRequest{
		APIKey:       apiKey,
		APIKeySecret: apiKeySecret,
		TargetUser:   TargetUser{UserID: userID},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("https://kasm.pkg.lat/api/public/get_login", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResp LoginResponse
	err = json.Unmarshal(respBytes, &loginResp)
	if err != nil {
		return "", err
	}

	return loginResp.URL, nil
}