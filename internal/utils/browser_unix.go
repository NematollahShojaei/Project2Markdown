//go:build !windows

// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"os/exec"
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
