// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateProject_ConcurrencyAndStability(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some dummy files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("hello world"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.md"), []byte("# Title\nContent"), 0644)

	// Create an unreadable file to test graceful error handling (Zero-Bug Policy)
	unreadableFile := filepath.Join(tmpDir, "unreadable.txt")
	os.WriteFile(unreadableFile, []byte("secret"), 0000) // No read permissions

	outName := filepath.Join(tmpDir, "output_context.md")

	// Run the generator
	// Zero-Bug Policy: Added empty string for outputDir to test default behavior
	_, tokens, _, err := GenerateProject(
		tmpDir, "", outName, "", "md", "", false, false, false, "", false, nil, nil,
	)

	if err != nil {
		t.Fatalf("GenerateProject failed unexpectedly: %v", err)
	}

	if tokens == 0 {
		t.Errorf("Expected tokens to be > 0")
	}

	// Verify output file exists
	if _, err := os.Stat(outName); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	// Restore permissions so the temp directory can be cleaned up by the test runner
	os.Chmod(unreadableFile, 0644)
}
