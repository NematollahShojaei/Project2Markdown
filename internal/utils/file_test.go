// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsTextFile(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Test Directory (Should return false safely)
	if IsTextFile(tmpDir) {
		t.Errorf("Expected directory to return false")
	}

	// 2. Test Empty File (Should return true)
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	os.WriteFile(emptyFile, []byte(""), 0644)
	if !IsTextFile(emptyFile) {
		t.Errorf("Expected empty file to return true")
	}

	// 3. Test Text File
	textFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(textFile, []byte("package main\n\nfunc main() {}"), 0644)
	if !IsTextFile(textFile) {
		t.Errorf("Expected text file to return true")
	}

	// 4. Test Binary File (Null bytes)
	binaryFile := filepath.Join(tmpDir, "test.bin")
	os.WriteFile(binaryFile, []byte{0x00, 0x01, 0x02, 0x03}, 0644)
	if IsTextFile(binaryFile) {
		t.Errorf("Expected binary file to return false")
	}
}
