//go:build !windows

// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

// OpenBrowser opens the default web browser on Unix-like systems (Linux, macOS).
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return nil
	}
	return cmd.Start()
}

// OpenFolderInOS opens the file manager on Unix-like systems.
func OpenFolderInOS(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		// xdg-open does not support file highlighting, so we open its parent directory
		dir := filepath.Dir(path)
		cmd = exec.Command("xdg-open", dir)
	case "darwin":
		// The "-R" flag in macOS reveals and highlights the file in Finder
		cmd = exec.Command("open", "-R", path)
	default:
		return nil
	}
	return cmd.Start()
}
