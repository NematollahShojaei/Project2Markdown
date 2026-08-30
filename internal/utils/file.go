// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

// IsTextFile safely reads the first 512 bytes of a file to determine if it is plain text.
// It ensures binary files (like .uasset, .png, .exe) are skipped automatically.
// Zero-Bug Policy: Handles empty files and edge cases gracefully.
func IsTextFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

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
