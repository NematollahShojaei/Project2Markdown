// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nematollahshojaei/project2markdown/internal/utils"
)

// JSONRestoreFormat defines the structure for restoring from a JSON context file.
type JSONRestoreFormat struct {
	Files []struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	} `json:"files"`
}

// RestoreProject reads a generated context file (MD, XML, JSON) and safely reconstructs the project structure.
// Security Pillar: Implements strict path sanitization to prevent Path Traversal attacks.
// Architecture Pillar: Decoupled from UI/CLI using the onProgress callback.
func RestoreProject(mdPath, destinationDir string, onProgress func(filesProcessed, currentTokens int, currentPath string)) error {
	var err error
	if destinationDir == "" {
		destinationDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get current directory: %v", err)
		}
	} else {
		destinationDir, err = filepath.Abs(destinationDir)
		if err != nil {
			return fmt.Errorf("invalid destination directory: %v", err)
		}
	}

	// Handle JSON Restore
	if strings.HasSuffix(strings.ToLower(mdPath), ".json") {
		return restoreFromJSON(mdPath, destinationDir, onProgress)
	}

	// Handle MD and XML Restore (Stream-based)
	file, err := os.Open(mdPath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, MaxLineBuffer)

	var currentFile *os.File
	// Zero-Bug Policy: Ensure the last opened file is safely closed to prevent descriptor leaks
	defer func() {
		if currentFile != nil {
			currentFile.Close()
		}
	}()
	inCodeBlock := false

	filesProcessed := 0
	totalTokens := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// 1. Detect File Header (MD or XML)
		rawPath := ""
		isMDHeader := false
		if strings.HasPrefix(line, "### File: ") {
			rawPath = strings.TrimSpace(strings.TrimPrefix(line, "### File: "))
			isMDHeader = true
		} else if strings.HasPrefix(line, "[FILE: ") && strings.HasSuffix(line, "]") {
			// Zero-Bug Policy: Legacy support for older generated files
			rawPath = strings.TrimSuffix(strings.TrimPrefix(line, "[FILE: "), "]")
			isMDHeader = true
		} else if strings.HasPrefix(trimmedLine, "<file path=\"") && strings.HasSuffix(trimmedLine, "\">") {
			rawPath = strings.TrimSuffix(strings.TrimPrefix(trimmedLine, "<file path=\""), "\">")
			isMDHeader = false
		}

		if rawPath != "" {
			cleanPath := filepath.Clean(rawPath)
			// Security Pillar: Prevent Path Traversal including Windows absolute paths (e.g., C:\)
			if strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
				fmt.Printf("Security Warning: Skipped malicious path: %s\n", cleanPath)
				continue
			}

			// Only strip the first directory if it's an MD file (which includes the ProjectName/ prefix)
			if isMDHeader {
				pathParts := strings.SplitN(filepath.ToSlash(cleanPath), "/", 2)
				if len(pathParts) == 2 {
					cleanPath = filepath.FromSlash(pathParts[1])
				}
			}

			absoluteTargetPath := filepath.Join(destinationDir, cleanPath)
			targetDir := filepath.Dir(absoluteTargetPath)
			os.MkdirAll(targetDir, 0755)

			if currentFile != nil {
				currentFile.Close()
			}
			currentFile, err = os.Create(absoluteTargetPath)
			if err == nil {
				filesProcessed++
				if onProgress != nil {
					onProgress(filesProcessed, totalTokens, cleanPath)
				}
			}
			continue
		}

		// 2. Handle Code Block Boundaries (MD or XML)
		if currentFile != nil {
			if strings.HasPrefix(line, "```") || strings.HasPrefix(trimmedLine, "<![CDATA[") {
				if !inCodeBlock {
					inCodeBlock = true
					continue
				}
			}
			if strings.HasPrefix(line, "```") || strings.HasPrefix(trimmedLine, "]]>") {
				if inCodeBlock {
					currentFile.Close()
					currentFile = nil
					inCodeBlock = false
					continue
				}
			}
		}

		// 3. Write Content to File
		if currentFile != nil && inCodeBlock {
			if strings.Contains(line, "` ` `") {
				line = strings.ReplaceAll(line, "` ` `", "```")
			}
			if strings.Contains(line, "]] >") {
				line = strings.ReplaceAll(line, "]] >", "]]>")
			}
			currentFile.WriteString(line + "\n")
			totalTokens += utils.EstimateTokens(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

// restoreFromJSON handles reverse generation specifically for JSON context files.
func restoreFromJSON(jsonPath, destinationDir string, onProgress func(filesProcessed, currentTokens int, currentPath string)) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var ctx JSONRestoreFormat
	if err := json.Unmarshal(data, &ctx); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	filesProcessed := 0
	totalTokens := 0

	for _, f := range ctx.Files {
		cleanPath := filepath.Clean(f.Path)
		// Security Pillar: Prevent Path Traversal including Windows absolute paths (e.g., C:\)
		if strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
			fmt.Printf("Security Warning: Skipped malicious path in JSON: %s\n", cleanPath)
			continue
		}

		absoluteTargetPath := filepath.Join(destinationDir, cleanPath)
		os.MkdirAll(filepath.Dir(absoluteTargetPath), 0755)

		err := os.WriteFile(absoluteTargetPath, []byte(f.Content), 0644)
		if err == nil {
			filesProcessed++
			totalTokens += utils.EstimateTokens(f.Content)
			if onProgress != nil {
				onProgress(filesProcessed, totalTokens, cleanPath)
			}
		}
	}

	return nil
}
