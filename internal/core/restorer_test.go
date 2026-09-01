// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRestoreProject_PathTraversalPrevention(t *testing.T) {
	tmpDir := t.TempDir()
	destDir := filepath.Join(tmpDir, "restored_project")
	os.MkdirAll(destDir, 0755)

	// Create a malicious Markdown context file
	maliciousMD := filepath.Join(tmpDir, "malicious.md")
	maliciousContent := `
### File: ProjectName/../../../malicious_outside.txt
` + "```text\nThis should not be written\n```\n" + `
### File: ProjectName/valid_inside.txt
` + "```text\nThis should be written\n```\n"

	err := os.WriteFile(maliciousMD, []byte(maliciousContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Run the restorer
	err = RestoreProject(maliciousMD, destDir, nil)
	if err != nil {
		t.Fatalf("RestoreProject failed: %v", err)
	}

	// Test Root Directory Edge Case
	err = RestoreProject(maliciousMD, "/", nil)
	if err != nil {
		t.Fatalf("RestoreProject failed on root directory: %v", err)
	}

	// 1. Verify the malicious file was NOT created outside the destination
	maliciousPath := filepath.Join(tmpDir, "malicious_outside.txt")
	if _, err := os.Stat(maliciousPath); !os.IsNotExist(err) {
		t.Errorf("SECURITY FAILURE: Path Traversal allowed! Malicious file was created at %s", maliciousPath)
	}

	// 2. Verify the valid file WAS created inside the destination
	validPath := filepath.Join(destDir, "valid_inside.txt")
	if _, err := os.Stat(validPath); os.IsNotExist(err) {
		t.Errorf("Expected valid file to be created at %s, but it was not", validPath)
	}
}
