package latex

import (
    "context"
    "fmt"
    "log"
   _ "os"
   _ "path/filepath"
    _"strings"
    "time"

    "google.golang.org/api/option"
    "github.com/google/generative-ai-go/genai"
)

const (
    maxGeminiAttempts = 4
    geminiTimeout     = 30 * time.Second
)

// FixLatexWithGemini attempts to fix LaTeX content using Gemini AI
func FixLatexWithGemini(apiKey, texContent, errorMsg string) (string, error) {
    if apiKey == "" {
        return "", fmt.Errorf("missing API key for Gemini")
    }

    ctx, cancel := context.WithTimeout(context.Background(), geminiTimeout)
    defer cancel()

    // Create Gemini client
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return "", fmt.Errorf("failed to create Gemini client: %w", err)
    }
    defer client.Close()

    // Use Gemini 2.5 Flash model
    model := client.GenerativeModel("models/gemini-2.5-pro")
    model.SetTemperature(0.2) // Lower temperature for more deterministic results
    model.SetMaxOutputTokens(2048)

    // Create prompt for Gemini
	prompt := fmt.Sprintf(`You are an expert LaTeX engineer whose sole job is to fix LaTeX sources so they compile with Tectonic. Using the ERROR MESSAGE and the LATEX DOCUMENT below, produce a corrected LaTeX source that will compile with Tectonic. Follow these rules strictly:
1) Diagnose the error from the provided message and make minimal, targeted fixes (syntax, missing braces, unclosed environments, incorrect environment names, missing math delimiters, mismatched \begin/\end, and missing common packages that are needed by the document).
2) Preserve the original document structure, macros, comments and intent; change only what is necessary to make it compile.
3) If adding packages is required, add only widely-available packages in the preamble (no external files). Prefer safety and compatibility with Tectonic.
4) Do not add explanations, diagnostics, or any text outside the LaTeX source. Do not use markdown or code fences.
5) If a best-effort fix still may have issues, return the best corrected LaTeX source you can produce (still with no explanations).
6) 

ERROR MESSAGE FOR THE TECTONIC LATEX ENGINE:
%s

LATEX DOCUMENT:
%s`, errorMsg, texContent)

    // Call Gemini API
    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return "", fmt.Errorf("Gemini API error: %w", err)
    }

    if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
        return "", fmt.Errorf("empty response from Gemini")
    }

    // Extract the response text
    responseText, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
    if !ok {
        return "", fmt.Errorf("unexpected response format from Gemini")
    }

    // Clean up the response - remove any markdown code block markers
    fixedLatex := RemoveCodeBlockMarkers(string(responseText))
    
    log.Printf("Successfully received fixed LaTeX from Gemini")
    return fixedLatex, nil
}
