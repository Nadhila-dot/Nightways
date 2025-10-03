package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

// GeminiRequest represents the request body for Gemini API
type GeminiRequest struct {
    Contents         []GeminiContent    `json:"contents"`
    SystemInstruction *GeminiInstruction `json:"systemInstruction,omitempty"`
}

// GeminiContent represents a content part
type GeminiContent struct {
    Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a text part
type GeminiPart struct {
    Text string `json:"text"`
}

// GeminiInstruction represents the system instruction
type GeminiInstruction struct {
    Parts []GeminiPart `json:"parts"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
    Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate response
type GeminiCandidate struct {
    Content GeminiContent `json:"content"`
}

// GenerateResponse generates a response using Gemini API
func GenerateResponse(apiKey, model, systemPrompt, userPrompt string, cooldownSec int) (string, error) {
    // Apply cooldown if specified
    if cooldownSec > 0 {
        time.Sleep(time.Duration(cooldownSec) * time.Second)
    }

    // Gemini API URL
    url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

    // Build request body
    reqBody := GeminiRequest{
        Contents: []GeminiContent{
            {
                Parts: []GeminiPart{
                    {Text: userPrompt},
                },
            },
        },
    }
    if systemPrompt != "" {
        reqBody.SystemInstruction = &GeminiInstruction{
            Parts: []GeminiPart{
                {Text: systemPrompt},
            },
        }
    }

    // Marshal to JSON
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %v", err)
    }

    // Make HTTP request
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to make request: %v", err)
    }
    defer resp.Body.Close()

    // Read response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API error: %s", string(body))
    }

    // Unmarshal response
    var geminiResp GeminiResponse
    if err := json.Unmarshal(body, &geminiResp); err != nil {
        return "", fmt.Errorf("failed to unmarshal response: %v", err)
    }

    // Extract text from response
    if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
        return geminiResp.Candidates[0].Content.Parts[0].Text, nil
    }

    return "", fmt.Errorf("no response generated")
}