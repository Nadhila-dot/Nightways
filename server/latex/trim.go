package latex

import (
    "fmt"
    "io/ioutil"
    "log"
   _ "os"
    "regexp"
    "strings"
)

// CleanTexFile removes Markdown code block markers from a LaTeX file
// and returns information about what was removed
func CleanTexFile(filePath string) (bool, string, error) {
    // Read the file
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return false, "", fmt.Errorf("failed to read file: %w", err)
    }
    
    originalContent := string(content)
    
    // Check if the content starts with ```latex or ```tex or ```
    hasMarkdownStart := false
    markdownPrefix := ""
    if strings.HasPrefix(originalContent, "```latex") {
        hasMarkdownStart = true
        markdownPrefix = "```latex"
    } else if strings.HasPrefix(originalContent, "```tex") {
        hasMarkdownStart = true
        markdownPrefix = "```tex"
    } else if strings.HasPrefix(originalContent, "```") {
        hasMarkdownStart = true
        markdownPrefix = "```"
    }
    
    // Check if the content ends with ```
    hasMarkdownEnd := false
    if strings.TrimSpace(originalContent)[len(strings.TrimSpace(originalContent))-3:] == "```" {
        hasMarkdownEnd = true
    }
    
    // If no markers found, return early
    if !hasMarkdownStart && !hasMarkdownEnd {
        return false, "No Markdown code block markers found.", nil
    }
    
    // Clean the content
    cleanedContent := originalContent
    
    // Remove leading markers
    if hasMarkdownStart {
        cleanedContent = strings.TrimPrefix(cleanedContent, markdownPrefix)
    }
    
    // Remove trailing markers
    if hasMarkdownEnd {
        cleanedContent = strings.TrimRight(cleanedContent, "`")
        // Make sure we don't trim too much - add back any needed newlines
        if !strings.HasSuffix(cleanedContent, "\n") {
            cleanedContent += "\n"
        }
    }
    
    // More aggressive removal of backtick markers with regex
    startPattern := regexp.MustCompile(`^` + "`{3,}" + `(latex|tex)?\s*(\r?\n)?`)
    cleanedContent = startPattern.ReplaceAllString(cleanedContent, "")
    
    endPattern := regexp.MustCompile(`\s*` + "`{3,}" + `\s*$`)
    cleanedContent = endPattern.ReplaceAllString(cleanedContent, "")
    
    // Trim whitespace but preserve leading/trailing newlines
    hasLeadingNewline := strings.HasPrefix(cleanedContent, "\n")
    hasTrailingNewline := strings.HasSuffix(cleanedContent, "\n")
    cleanedContent = strings.TrimSpace(cleanedContent)
    if hasLeadingNewline {
        cleanedContent = "\n" + cleanedContent
    }
    if hasTrailingNewline {
        cleanedContent = cleanedContent + "\n"
    }
    
    // Create summary of changes
    summary := fmt.Sprintf("Removed Markdown markers:\n")
    if hasMarkdownStart {
        summary += fmt.Sprintf("- Leading marker: %s\n", markdownPrefix)
    }
    if hasMarkdownEnd {
        summary += fmt.Sprintf("- Trailing marker: ```\n")
    }
    
    // Write the cleaned content back to the file
    err = ioutil.WriteFile(filePath, []byte(cleanedContent), 0644)
    if err != nil {
        return false, "", fmt.Errorf("failed to write cleaned content: %w", err)
    }
    
    log.Printf("Successfully cleaned LaTeX file: %s", filePath)
    return true, summary, nil
}

