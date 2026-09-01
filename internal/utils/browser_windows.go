//go:build windows

// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import (
	"os/exec"
	"path/filepath"
	"syscall"
)

// OpenBrowser opens the default web browser without flashing a CMD window on Windows.
func OpenBrowser(url string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)

	// Zero-Bug Policy: Prevent Windows from spawning a temporary CMD window for the browser
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Start()
}

// OpenFolderInOS opens the Windows File Explorer and highlights the generated file.
func OpenFolderInOS(path string) error {
	cleanPath := filepath.FromSlash(path)

	// Zero-Bug Policy: Call explorer directly without cmd /c to prevent any black console windows.
	// Explorer is a GUI subsystem app, so it won't flash a console natively.
	cmd := exec.Command("explorer", "/select,"+cleanPath)

	return cmd.Start()
}
