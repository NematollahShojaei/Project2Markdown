//go:build windows

// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"os/exec"
	"syscall"
)

// OpenBrowser opens the default web browser without flashing a CMD window on Windows.
func OpenBrowser(url string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)

	// Zero-Bug Policy: Prevent Windows from spawning a temporary CMD window for the browser
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Start()
}
