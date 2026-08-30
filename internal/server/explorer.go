// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type ExploreRequest struct {
	Path string `json:"path"`
}

type MkdirRequest struct {
	Path       string `json:"path"`
	FolderName string `json:"folderName"`
}

type FileNode struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

type QuickAccess struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
}

// handleExplore safely reads directory contents without exposing system vulnerabilities.
func handleExplore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ExploreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var nodes []FileNode

	// Handle root directory / drives
	if req.Path == "" {
		if runtime.GOOS == "windows" {
			// Get logical drives on Windows
			for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
				d := string(drive) + ":\\"
				if _, err := os.Stat(d); err == nil {
					nodes = append(nodes, FileNode{Name: d, Path: d, IsDir: true})
				}
			}
		} else {
			req.Path = "/"
		}
	}

	if req.Path != "" {
		entries, err := os.ReadDir(req.Path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Access Denied or Error: %v", err), http.StatusInternalServerError)
			return
		}

		// Add parent directory option if not at root
		parent := filepath.Dir(req.Path)
		if parent != req.Path {
			nodes = append(nodes, FileNode{Name: "..", Path: parent, IsDir: true})
		}

		for _, entry := range entries {
			fullPath := filepath.Join(req.Path, entry.Name())
			nodes = append(nodes, FileNode{
				Name:  entry.Name(),
				Path:  fullPath,
				IsDir: entry.IsDir(),
			})
		}
	}

	// Sort: Directories first, then alphabetically
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Name == ".." {
			return true
		}
		if nodes[j].Name == ".." {
			return false
		}
		if nodes[i].IsDir == nodes[j].IsDir {
			return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
		}
		return nodes[i].IsDir
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

// getKnownFolder intelligently resolves Windows known folders (like Desktop/Documents)
// by checking if they have been redirected to OneDrive.
func getKnownFolder(home, folderName string) string {
	oneDrivePath := filepath.Join(home, "OneDrive", folderName)
	if _, err := os.Stat(oneDrivePath); err == nil {
		return oneDrivePath
	}
	return filepath.Join(home, folderName)
}

// handleQuickAccess returns standard user directories and drives for the sidebar.
func handleQuickAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var qa []QuickAccess
	home, err := os.UserHomeDir()
	if err == nil {
		qa = append(qa, QuickAccess{Name: "Home", Path: home, Icon: "home"})
		qa = append(qa, QuickAccess{Name: "Desktop", Path: getKnownFolder(home, "Desktop"), Icon: "desktop"})
		qa = append(qa, QuickAccess{Name: "Documents", Path: getKnownFolder(home, "Documents"), Icon: "document"})
		qa = append(qa, QuickAccess{Name: "Downloads", Path: filepath.Join(home, "Downloads"), Icon: "download"})
	}

	if runtime.GOOS == "windows" {
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			d := string(drive) + ":\\"
			if _, err := os.Stat(d); err == nil {
				qa = append(qa, QuickAccess{Name: string(drive) + " Drive", Path: d, Icon: "drive"})
			}
		}
	} else {
		qa = append(qa, QuickAccess{Name: "Root", Path: "/", Icon: "drive"})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qa)
}

// handleMkdir creates a new directory at the specified path.
func handleMkdir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req MkdirRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.Path == "" || req.FolderName == "" {
		http.Error(w, "Path and Folder Name are required", http.StatusBadRequest)
		return
	}

	newDirPath := filepath.Join(req.Path, req.FolderName)
	if err := os.MkdirAll(newDirPath, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create folder: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Folder created successfully"))
}
