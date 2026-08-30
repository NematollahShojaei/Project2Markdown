// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package remote

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// progressWriter tracks the download progress and triggers a callback.
type progressWriter struct {
	writer     io.Writer
	total      int64
	downloaded int64
	onProgress func(downloaded, total int64)
	lastReport time.Time
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.downloaded += int64(n)

	// Throttle updates to every 500ms to prevent UI lag
	if pw.onProgress != nil && time.Since(pw.lastReport) > 500*time.Millisecond {
		pw.onProgress(pw.downloaded, pw.total)
		pw.lastReport = time.Now()
	}
	return n, err
}

// DownloadAndExtractRepo downloads a GitHub repository zipball and extracts it safely.
// Returns the path to the extracted source, the temp directory to clean up, and any error.
func DownloadAndExtractRepo(repoURL string, onProgress func(downloaded, total int64)) (string, string, error) {
	repoURL = strings.TrimPrefix(repoURL, "https://github.com/")
	repoURL = strings.TrimPrefix(repoURL, "http://github.com/")
	repoURL = strings.TrimPrefix(repoURL, "github.com/")
	parts := strings.Split(strings.Trim(repoURL, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub repository format. Use owner/repo")
	}
	owner, repo := parts[0], parts[1]

	// Security Pillar: Sanitize GitHub owner and repo to prevent URL injection/SSRF
	if strings.Contains(owner, ".") || strings.Contains(owner, "/") || strings.Contains(repo, ".") || strings.Contains(repo, "/") {
		return "", "", fmt.Errorf("invalid GitHub repository name format")
	}

	zipURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball", owner, repo)

	tmpZip, err := os.CreateTemp("", "p2m-*.zip")
	if err != nil {
		return "", "", err
	}
	zipPath := tmpZip.Name()
	tmpZip.Close()

	// Performance & Stability Pillar: Use a custom HTTP client with a strict timeout to prevent infinite hangs
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(zipURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to download repository: HTTP %d (Check if repo is public)", resp.StatusCode)
	}

	out, err := os.Create(zipPath)
	if err != nil {
		return "", "", err
	}

	pw := &progressWriter{
		writer:     out,
		total:      resp.ContentLength, // Can be -1 if the server doesn't send Content-Length
		onProgress: onProgress,
		lastReport: time.Now(),
	}

	_, err = io.Copy(pw, resp.Body)
	out.Close()
	if err != nil {
		return "", "", err
	}

	tmpDir, err := os.MkdirTemp("", "p2m-repo-*")
	if err != nil {
		return "", "", err
	}

	extractedRoot, err := unzipSafe(zipPath, tmpDir)
	os.Remove(zipPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", "", err
	}

	if extractedRoot == "" {
		extractedRoot = tmpDir
	}
	return extractedRoot, tmpDir, nil
}

// unzipSafe extracts a zip archive with strict Zip Slip vulnerability prevention.
func unzipSafe(src, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var rootDir string
	for i, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Security Pillar: Prevent Zip Slip attacks by ensuring the file path is within the destination directory
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return "", fmt.Errorf("illegal file path (Zip Slip detected): %s", fpath)
		}

		if i == 0 || rootDir == "" {
			parts := strings.Split(f.Name, "/")
			if len(parts) > 0 {
				rootDir = filepath.Join(dest, parts[0])
			}
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return "", err
		}
	}
	return rootDir, nil
}
