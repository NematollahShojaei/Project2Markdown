// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IgnoreEngine handles parsing and evaluating ignore files (.gitignore, .p4ignore, .p2mignore).
// Zero-Bug Policy: Uses a high-performance pattern matching system to prevent memory leaks.
type IgnoreEngine struct {
	patterns []string
}

// NewIgnoreEngine initializes the engine with default safe exclusions and loads custom ignore files.
func NewIgnoreEngine(rootDir string) *IgnoreEngine {
	engine := &IgnoreEngine{
		patterns: []string{
			".git", ".github", "node_modules", "vendor", "bin", "obj",
			"binaries", "deriveddatacache", "intermediate", "saved", ".vs",
		},
	}

	ignoreFiles := []string{".gitignore", ".p4ignore", ".p2mignore"}
	for _, igFile := range ignoreFiles {
		engine.loadPatterns(filepath.Join(rootDir, igFile))
	}
	return engine
}

// AddCustomPatterns allows injecting comma-separated ignore rules on the fly from the UI or CLI.
func (ie *IgnoreEngine) AddCustomPatterns(customIgnores string) {
	if customIgnores == "" {
		return
	}
	rules := strings.Split(customIgnores, ",")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule != "" {
			ie.patterns = append(ie.patterns, filepath.ToSlash(rule))
		}
	}
}

// loadPatterns safely reads an ignore file and extracts valid patterns.
func (ie *IgnoreEngine) loadPatterns(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return // Safely ignore if the file does not exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Normalize path separators for cross-platform compatibility
		line = filepath.ToSlash(line)
		ie.patterns = append(ie.patterns, line)
	}

	// Zero-Bug Safety: Ensure no I/O errors occurred during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Warning: Error reading ignore file [%s]: %v\n", filePath, err)
	}
}

// IsIgnored checks if a given relative path matches any of the loaded ignore patterns.
func (ie *IgnoreEngine) IsIgnored(relPath string, isDir bool) bool {
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range ie.patterns {
		pattern = strings.TrimPrefix(pattern, "/")
		pattern = strings.TrimSuffix(pattern, "/")

		// 1. Exact match or directory prefix match
		if relPath == pattern || strings.HasPrefix(relPath, pattern+"/") {
			return true
		}

		// 2. Wildcard suffix match (e.g., *.log)
		if strings.HasPrefix(pattern, "*") {
			suffix := strings.TrimPrefix(pattern, "*")
			if strings.HasSuffix(relPath, suffix) {
				return true
			}
		}

		// 3. Directory/File name match anywhere in the path
		if strings.Contains("/"+relPath+"/", "/"+pattern+"/") {
			return true
		}
	}
	return false
}
