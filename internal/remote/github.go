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
	"regexp"
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

	// Security Pillar: Strict Regex validation to prevent SSRF and URL Injection
	// GitHub usernames and repos must only contain alphanumeric characters, hyphens, underscores, or periods.
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !validNameRegex.MatchString(owner) || !validNameRegex.MatchString(repo) {
		return "", "", fmt.Errorf("invalid GitHub repository name format (contains illegal characters)")
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

	// Zero-Bug Policy: Wrap file operations in an anonymous function to ensure defer executes immediately after copy
	err = func() error {
		out, err := os.Create(zipPath)
		if err != nil {
			return err
		}
		defer out.Close()

		pw := &progressWriter{
			writer:     out,
			total:      resp.ContentLength,
			onProgress: onProgress,
			lastReport: time.Now(),
		}

		_, err = io.Copy(pw, resp.Body)
		return err
	}()

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
	destClean := filepath.Clean(dest) + string(os.PathSeparator)

	for i, f := range r.File {
		// Zero-Bug Policy: Use closure to guarantee defer execution per iteration (prevents FD exhaustion)
		err := func() error {
			fpath := filepath.Join(dest, f.Name)

			// Security Pillar: Prevent Zip Slip attacks by ensuring the file path is strictly within the destination directory
			if !strings.HasPrefix(fpath, destClean) {
				return fmt.Errorf("illegal file path (Zip Slip detected): %s", fpath)
			}

			if i == 0 || rootDir == "" {
				parts := strings.Split(f.Name, "/")
				if len(parts) > 0 {
					rootDir = filepath.Join(dest, parts[0])
				}
			}

			if f.FileInfo().IsDir() {
				return os.MkdirAll(fpath, os.ModePerm)
			}

			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			_, err = io.Copy(outFile, rc)
			return err
		}()

		if err != nil {
			return "", err
		}
	}
	return rootDir, nil
}
