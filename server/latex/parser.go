package latex

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ParseResult contains information about parsed LaTeX content
type ParseResult struct {
	IsValid      bool              `json:"isValid"`
	HasErrors    bool              `json:"hasErrors"`
	Warnings     []string          `json:"warnings"`
	Errors       []string          `json:"errors"`
	Stats        map[string]int    `json:"stats"`
	Structure    map[string]string `json:"structure"`
	Dependencies []string          `json:"dependencies"`
}

// ParseLaTeX analyzes LaTeX content for validity and structure
func ParseLaTeX(content string) (*ParseResult, error) {
	fmt.Println(content) // debug content

	if content == "" {
		return nil, fmt.Errorf("empty LaTeX content")
	}

	result := &ParseResult{
		IsValid:      true,
		HasErrors:    false,
		Warnings:     []string{},
		Errors:       []string{},
		Stats:        make(map[string]int),
		Structure:    make(map[string]string),
		Dependencies: []string{},
	}

	// Check for document class
	documentClassRegex := regexp.MustCompile(`\\documentclass(\[.*?\])?\{(.*?)\}`)
	documentClassMatch := documentClassRegex.FindStringSubmatch(content)

	if len(documentClassMatch) < 3 {
		result.IsValid = false
		result.HasErrors = true
		result.Errors = append(result.Errors, "Missing \\documentclass declaration")
	} else {
		result.Structure["documentClass"] = documentClassMatch[2]
	}

	// Check for document environment
	if !strings.Contains(content, "\\begin{document}") || !strings.Contains(content, "\\end{document}") {
		result.IsValid = false
		result.HasErrors = true
		result.Errors = append(result.Errors, "Missing document environment")
	}

	// Extract packages
	packageRegex := regexp.MustCompile(`\\usepackage(\[.*?\])?\{(.*?)\}`)
	packageMatches := packageRegex.FindAllStringSubmatch(content, -1)

	packages := []string{}
	for _, match := range packageMatches {
		if len(match) >= 3 {
			packages = append(packages, match[2])
		}
	}
	result.Dependencies = packages
	result.Stats["packageCount"] = len(packages)

	// Count sections
	sectionRegex := regexp.MustCompile(`\\section\{`)
	sectionMatches := sectionRegex.FindAllStringIndex(content, -1)
	result.Stats["sectionCount"] = len(sectionMatches)

	// Count subsections
	subsectionRegex := regexp.MustCompile(`\\subsection\{`)
	subsectionMatches := subsectionRegex.FindAllStringIndex(content, -1)
	result.Stats["subsectionCount"] = len(subsectionMatches)

	// Count environments
	environmentRegex := regexp.MustCompile(`\\begin\{(.*?)\}`)
	environmentMatches := environmentRegex.FindAllStringSubmatch(content, -1)

	environments := make(map[string]int)
	for _, match := range environmentMatches {
		if len(match) >= 2 {
			env := match[1]
			environments[env]++
		}
	}

	// Check for environment balance
	for env, count := range environments {
		beginCount := count
		endCount := strings.Count(content, "\\end{"+env+"}")

		if beginCount != endCount {
			result.HasErrors = true
			result.Errors = append(result.Errors,
				fmt.Sprintf("Unbalanced environment: %s (begin: %d, end: %d)",
					env, beginCount, endCount))
		}
	}

	// Count math environments
	mathEnvs := []string{"equation", "align", "math", "displaymath"}
	mathCount := 0
	for _, env := range mathEnvs {
		if count, ok := environments[env]; ok {
			mathCount += count
		}
	}
	// Also count inline math
	inlineMathRegex := regexp.MustCompile(`\$[^\$]+\$`)
	inlineMathMatches := inlineMathRegex.FindAllStringIndex(content, -1)
	mathCount += len(inlineMathMatches)

	result.Stats["mathCount"] = mathCount

	// Check for common warnings
	if !strings.Contains(content, "\\title{") {
		result.Warnings = append(result.Warnings, "No title defined")
	}

	RemoveCodeBlockMarkers(content) // Clean content from code block markers

	return result, nil
}

// ExtractOutput extracts LaTeX content from a string between <Output> tags
func ExtractOutput(content string) (string, error) {
	outputRegex := regexp.MustCompile(`(?s)<Output>(.*?)</Output>`)
	matches := outputRegex.FindStringSubmatch(content)

	if len(matches) < 2 {
		return "", fmt.Errorf("could not find output content")
	}

	return strings.TrimSpace(matches[1]), nil
}

// ExtractMetadata extracts metadata from a string between <meta-data> tags
func ExtractMetadata(content string) (map[string]string, error) {
	metadataRegex := regexp.MustCompile(`(?s)<meta-data>(.*?)</meta-data>`)
	matches := metadataRegex.FindStringSubmatch(content)

	if len(matches) < 2 {
		return nil, fmt.Errorf("could not find metadata content")
	}

	metadata := make(map[string]string)
	metadataText := matches[1]
	lines := strings.Split(metadataText, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			metadata[key] = value
		}
	}

	return metadata, nil
}

// SplitContent takes raw AI output and splits it into LaTeX and metadata components
func SplitContent(content string) (string, map[string]string, error) {
	// Clean any code block markers first
	content = RemoveCodeBlockMarkers(content)

	// Check if content has Output tags
	outputRegex := regexp.MustCompile(`(?s)<Output>(.*?)</Output>`)
	outputMatches := outputRegex.FindStringSubmatch(content)

	var latex string
	if len(outputMatches) < 2 {
		// No Output tags found, check if the content looks like LaTeX code
		if strings.Contains(content, "\\documentclass") ||
			(strings.Contains(content, "\\begin{document}") &&
				strings.Contains(content, "\\end{document}")) {
			// Content appears to be direct LaTeX
			latex = content // Already cleaned above
		} else if strings.Contains(content, "\\") &&
			(strings.Contains(content, "{") || strings.Contains(content, "}")) {
			// More relaxed check - any content with LaTeX commands is probably LaTeX
			latex = content
		} else {
			// Return the content as-is as a last resort rather than failing
			latex = content
			return latex, map[string]string{"source": "raw-content", "generated": time.Now().Format(time.RFC3339)}, nil
		}
	} else {
		// Extract content from Output tags
		latex = strings.TrimSpace(outputMatches[1])
	}

	// Extract metadata if available, or return empty metadata
	metadata := make(map[string]string)
	metadataRegex := regexp.MustCompile(`(?s)<meta-data>(.*?)</meta-data>`)
	metadataMatches := metadataRegex.FindStringSubmatch(content)

	if len(metadataMatches) >= 2 {
		metadataText := metadataMatches[1]
		lines := strings.Split(metadataText, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				metadata[key] = value
			}
		}
	} else {
		// If no metadata found, add some basic info
		metadata["source"] = "direct-latex"
		metadata["generated"] = time.Now().Format(time.RFC3339)
	}

	return latex, metadata, nil
}

func RemoveCodeBlockMarkers(content string) string {
    // Keep a copy of the original content
    originalContent := content

    // Trim whitespace before processing
    content = strings.TrimSpace(content)

    // Remove explicit ```latex marker at the beginning
    if strings.HasPrefix(content, "```latex") {
        content = strings.TrimPrefix(content, "```latex")
    } else if strings.HasPrefix(content, "```tex") {
        content = strings.TrimPrefix(content, "```tex")
    } else if strings.HasPrefix(content, "```") {
        content = strings.TrimPrefix(content, "```")
    }

    // Remove closing backticks
    if strings.HasSuffix(content, "```") {
        content = strings.TrimSuffix(content, "```")
    }

    // More aggressive removal of backtick markers with regex
    // Match any number of backticks (3 or more) followed by optional language identifier
    startPattern := regexp.MustCompile(`^` + "`{3,}" + `(latex|tex)?\s*(\r?\n)?`)
    content = startPattern.ReplaceAllString(content, "")

    // Match any number of backticks at the end
    endPattern := regexp.MustCompile(`\s*` + "`{3,}" + `\s*$`)
    content = endPattern.ReplaceAllString(content, "")

    // Trim whitespace again
    content = strings.TrimSpace(content)

    // Basic validation - if we ended up with nothing useful, return the original
    if len(content) < 20 || !strings.Contains(content, "\\documentclass") {
        log.Printf("Content cleaning removed too much or important markers, reverting to original")
        return originalContent
    }

    log.Printf("Successfully removed code block markers from LaTeX content")
    return content
}

// GenerationResult holds generated LaTeX and its metadata
type GenerationResult struct {
	LaTeX   string
	Metadata map[string]string
}

// SaveGeneratedContent saves the LaTeX content and metadata to separate files
func SaveGeneratedContent(outputDir, jobID string, result *GenerationResult) error {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save original LaTeX content to beforeclean directory
	beforeCleanDir := filepath.Join(outputDir, "beforeclean")
	if err := os.MkdirAll(beforeCleanDir, 0755); err != nil {
		return fmt.Errorf("failed to create beforeclean directory: %w", err)
	}
	beforeCleanFile := filepath.Join(beforeCleanDir, jobID+".tex")
	if err := os.WriteFile(beforeCleanFile, []byte(result.LaTeX), 0644); err != nil {
		return fmt.Errorf("failed to save beforeclean LaTeX file: %w", err)
	}
	log.Printf("Successfully saved beforeclean .tex file to: %s", beforeCleanFile)

	// Clean the LaTeX content to ensure no markdown code markers exist
	cleanedLatex := RemoveCodeBlockMarkers(result.LaTeX)
	cleanedLatex = cleanLatexContent(cleanedLatex)

	// Save cleaned LaTeX content to .tex file

	texFile := filepath.Join(outputDir, jobID+".tex")
	// Only remove starting ```latex/```tex/``` and ending ``` markers, not inside LaTeX content
	startPattern := regexp.MustCompile(`^` + "`{3,}" + `(latex|tex)?\s*(\r?\n)?`)
	cleanedLatex = startPattern.ReplaceAllString(cleanedLatex, "")
	endPattern := regexp.MustCompile(`\s*` + "`{3,}" + `\s*$`)
	cleanedLatex = endPattern.ReplaceAllString(cleanedLatex, "")

	err := os.WriteFile(texFile, []byte(cleanedLatex), 0644)
	if err != nil {
		return fmt.Errorf("failed to save LaTeX file: %w", err)
	}
	log.Printf("Successfully saved .tex file to: %s", texFile)

	// Save metadata as JSON
	metaFile := filepath.Join(outputDir, jobID+".meta.json")
	metaJSON, err := json.MarshalIndent(result.Metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = os.WriteFile(metaFile, metaJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to save metadata file: %w", err)
	}
	log.Printf("Successfully saved metadata to: %s", metaFile)

	return nil
}

// cleanLatexContent ensures LaTeX content is free of markdown code block markers
func cleanLatexContent(content string) string {
	// Trim whitespace
	content = strings.TrimSpace(content)

	// Remove common markdown code block markers
	if strings.HasPrefix(content, "```latex") {
		content = strings.TrimPrefix(content, "```latex")
	} else if strings.HasPrefix(content, "```tex") {
		content = strings.TrimPrefix(content, "```tex")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}

	// Handle the closing backticks
	if strings.HasSuffix(content, "```") {
		content = strings.TrimSuffix(content, "```")
	}

	// Remove any additional backtick markers that may be at the beginning or end
	backtickPattern := regexp.MustCompile(`^` + "`{3,}" + `(latex|tex)?\s*`)
	content = backtickPattern.ReplaceAllString(content, "")

	endBacktickPattern := regexp.MustCompile("`{3,}\\s*$")
	content = endBacktickPattern.ReplaceAllString(content, "")

	return strings.TrimSpace(content)
}
