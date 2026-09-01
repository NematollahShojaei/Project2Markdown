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
	patterns         []string
	blacklist        []string
	explicitIncludes []string
	allowSecrets     bool
}

// NewIgnoreEngine initializes the engine with default safe exclusions and loads custom ignore files.
func NewIgnoreEngine(rootDir string) *IgnoreEngine {
	engine := &IgnoreEngine{
		patterns: []string{
			".git", ".github", "node_modules", "vendor", "bin", "obj",
			"binaries", "deriveddatacache", "intermediate", "saved", ".vs",
		},
		// Security Pillar: Hardcoded blacklist for sensitive files to prevent accidental leaks
		blacklist: []string{
			".env", "*.env.*", "*.pem", "*.key", "*.p12", "*.pfx",
			"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519", "secrets.json", "*.keystore",
		},
		allowSecrets: false,
	}

	ignoreFiles := []string{".gitignore", ".p4ignore", ".p2mignore"}
	for _, igFile := range ignoreFiles {
		engine.loadPatterns(filepath.Join(rootDir, igFile))
	}
	return engine
}

// SetAllowSecrets toggles the security blacklist.
func (ie *IgnoreEngine) SetAllowSecrets(allow bool) {
	ie.allowSecrets = allow
}

// SetExplicitIncludes registers patterns that should bypass all ignore rules.
func (ie *IgnoreEngine) SetExplicitIncludes(includes string) {
	if includes == "" {
		return
	}
	rules := strings.Split(includes, ",")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule != "" {
			ie.explicitIncludes = append(ie.explicitIncludes, filepath.ToSlash(rule))
		}
	}
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

// IsIgnored checks if a given relative path matches standard ignore patterns.
func (ie *IgnoreEngine) IsIgnored(relPath string, isDir bool) bool {
	relPath = filepath.ToSlash(relPath)

	for _, inc := range ie.explicitIncludes {
		if strings.Contains(relPath, inc) || strings.HasSuffix(relPath, inc) {
			return false
		}
	}

	for _, pattern := range ie.patterns {
		if matchPattern(relPath, pattern) {
			return true
		}
	}
	return false
}

// IsSecurityBlocked checks if a file is sensitive and its content should be redacted.
func (ie *IgnoreEngine) IsSecurityBlocked(relPath string) bool {
	relPath = filepath.ToSlash(relPath)

	if ie.allowSecrets {
		return false
	}

	for _, inc := range ie.explicitIncludes {
		if strings.Contains(relPath, inc) || strings.HasSuffix(relPath, inc) {
			return false
		}
	}

	for _, b := range ie.blacklist {
		if matchPattern(relPath, b) {
			return true
		}
	}
	return false
}

// matchPattern is a helper function to evaluate path matches safely.
func matchPattern(relPath, pattern string) bool {
	pattern = strings.TrimPrefix(pattern, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	// Exact match or directory prefix match
	if relPath == pattern || strings.HasPrefix(relPath, pattern+"/") {
		return true
	}

	// Wildcard suffix match (e.g., *.log)
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		if strings.HasSuffix(relPath, suffix) {
			return true
		}
	}

	// Directory/File name match anywhere in the path
	if strings.Contains("/"+relPath+"/", "/"+pattern+"/") {
		return true
	}

	return false
}
