// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// IsTextFile safely reads the first 512 bytes of a file to determine if it is plain text.
// It ensures binary files (like .uasset, .png, .exe) are skipped automatically.
// Zero-Bug Policy: Handles empty files, directories, and permission edge cases gracefully.
func IsTextFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false // Safely skip unreadable files (e.g., permission denied)
	}
	defer file.Close()

	// Security Pillar: Ensure the path is strictly a regular file.
	// Zero-Bug Policy: Reject Named Pipes (FIFO), Sockets, and Devices which can cause infinite blocking on Read().
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return false
	}

	// Read up to 512 bytes for content type detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// Empty files are considered text (safe to include, though empty)
	if n == 0 {
		return true
	}

	// Use standard HTTP content type detection
	contentType := http.DetectContentType(buffer[:n])

	// http.DetectContentType returns "text/plain; charset=utf-8" for text files.
	if strings.HasPrefix(contentType, "text/") {
		return true
	}

	// Fallback: Check for null bytes (\x00).
	// If no null bytes exist in the first 512 bytes, it is highly likely a safe text file (e.g., JSON, XML, Dockerfile).
	if bytes.IndexByte(buffer[:n], 0) == -1 {
		return true
	}

	return false
}

// IsSafeCommentLine determines if a line is purely a comment.
// Zero-Bug Policy: Only matches lines that start with comment tokens to avoid breaking URLs or strings.
func IsSafeCommentLine(line, ext string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "*/") || strings.HasPrefix(trimmed, "<!--") {
		return true
	}

	// Safe handling for '#' comments (Python, Ruby, Shell, YAML, etc.)
	if strings.HasPrefix(trimmed, "#") {
		if ext == ".py" || ext == ".rb" || ext == ".sh" || ext == ".yaml" || ext == ".yml" || ext == ".dockerfile" || ext == ".env" || ext == "" {
			return true
		}
	}
	return false
}

// EscapeJSONString safely escapes a string for JSON embedding without heavy libraries.
func EscapeJSONString(s string) string {
	b, _ := json.Marshal(s)
	str := string(b)
	if len(str) >= 2 {
		return str[1 : len(str)-1]
	}
	return str
}

// IsWritable checks if the application has write permissions for the given directory.
// Zero-Bug Policy: The most reliable cross-platform method in Go is attempting to create a temporary file.
func IsWritable(dir string) bool {
	tempFile := filepath.Join(dir, ".p2m_write_test")
	f, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(tempFile)
	return true
}

// GetDefaultWorkspace returns the path to Documents/Project2Markdown.
// It handles OneDrive redirection on Windows automatically.
func GetDefaultWorkspace() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	docsPath := filepath.Join(home, "Documents")

	// Check for OneDrive redirection (Windows)
	oneDriveDocs := filepath.Join(home, "OneDrive", "Documents")
	if stat, err := os.Stat(oneDriveDocs); err == nil && stat.IsDir() {
		docsPath = oneDriveDocs
	}

	workspace := filepath.Join(docsPath, "Project2Markdown")
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return "", err
	}

	return workspace, nil
}
