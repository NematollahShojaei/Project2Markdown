// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package remote

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadAndExtractRepo_InvalidURLs(t *testing.T) {
	// Test SSRF and URL Injection prevention
	invalidURLs := []string{
		"owner/repo/extra",
		"owner_with_invalid_chars!/repo",
		"owner/repo_with_invalid_chars@",
		"../owner/repo",
		"owner/../../repo",
	}

	for _, url := range invalidURLs {
		_, _, err := DownloadAndExtractRepo(url, nil)
		if err == nil {
			t.Errorf("Expected error for invalid URL: %s, but got none", url)
		}
	}
}

func TestUnzipSafe_ZipSlipPrevention(t *testing.T) {
	// Create a malicious zip file in memory
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Add a file with a path traversal payload
	f, err := w.Create("../malicious.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte("malicious content"))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	// Write the buffer to a temporary file
	tmpZip := filepath.Join(t.TempDir(), "malicious.zip")
	err = os.WriteFile(tmpZip, buf.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to extract it
	destDir := t.TempDir()
	_, err = unzipSafe(tmpZip, destDir)

	// It MUST fail with a Zip Slip error
	if err == nil {
		t.Fatal("Expected unzipSafe to fail due to Zip Slip vulnerability, but it succeeded")
	}
}
