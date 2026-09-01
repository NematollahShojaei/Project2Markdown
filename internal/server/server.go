// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nematollahshojaei/project2markdown/v2/internal/core"
	"github.com/nematollahshojaei/project2markdown/v2/internal/remote"
	"github.com/nematollahshojaei/project2markdown/v2/internal/utils"
)

// API Request Structs for JSON parsing
type GenerateRequest struct {
	TargetDir        string `json:"targetDir"`
	OutputDir        string `json:"outputDir"`
	RemoteURL        string `json:"remoteURL"`
	OutputFileName   string `json:"outputFileName"`
	CustomIgnores    string `json:"customIgnores"`
	Format           string `json:"format"`
	IncludePatterns  string `json:"includePatterns"`
	RemoveComments   bool   `json:"removeComments"`
	RemoveEmptyLines bool   `json:"removeEmptyLines"`
	AiInstructions   bool   `json:"aiInstructions"`
	CustomPrompt     string `json:"customPrompt"`
	AllowSecrets     bool   `json:"allowSecrets"`
}

type RestoreRequest struct {
	MdFilePath     string `json:"mdFilePath"`
	DestinationDir string `json:"destinationDir"`
}

// StartUI initializes the local web server, registers API routes, and opens the browser.
// Architecture Pillar: Accepts uiHTML as a dependency to remain decoupled from the file system.
func StartUI(port string, uiHTML []byte) {
	fmt.Printf("Starting P2M Web UI on http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server.")

	// Initialize the SSE Broker locally to avoid Global State
	broker := NewBroker()

	var (
		lastHeartbeat time.Time
		heartbeatMu   sync.Mutex
		shutdownTimer *time.Timer
	)

	// Route: Serve the embedded HTML UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(uiHTML)
	})

	// Route: Generate Context
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req GenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		targetDir := req.TargetDir
		outName := req.OutputFileName
		outputDir := req.OutputDir
		cleanupDir := ""

		// Handle Remote GitHub Repository
		if req.RemoteURL != "" {
			onDownloadProgress := func(downloaded, total int64) {
				mb := float64(downloaded) / (1024 * 1024)
				if total > 0 {
					totalMB := float64(total) / (1024 * 1024)
					broker.SendLog("INFO", "Downloading", fmt.Sprintf("%.1f MB / %.1f MB", mb, totalMB))
				} else {
					broker.SendLog("INFO", "Downloading", fmt.Sprintf("%.1f MB", mb))
				}
			}

			extractedPath, tmpDir, err := remote.DownloadAndExtractRepo(req.RemoteURL, onDownloadProgress)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error downloading repo: %v", err), http.StatusInternalServerError)
				return
			}
			targetDir = extractedPath
			cleanupDir = tmpDir

			// Zero-Bug Policy: Prevent saving the output file inside the temporary directory
			if outputDir == "" {
				workspace, err := utils.GetDefaultWorkspace()
				if err == nil {
					outputDir = workspace
				}
			}

			if outName == "" {
				parts := strings.Split(strings.Trim(strings.TrimPrefix(strings.TrimPrefix(req.RemoteURL, "https://github.com/"), "github.com/"), "/"), "/")
				repoName := "repo"
				if len(parts) >= 2 {
					repoName = parts[1]
				}
				outName = repoName + "_context." + req.Format
			}
		}

		startTime := time.Now()
		totalFilesProcessed := 0

		// Progress Callback for Real-time SSE Updates
		onProgress := func(filesProcessed, currentTokens int, currentPath string) {
			totalFilesProcessed = filesProcessed
			broker.SendLog("INFO", "Processing", currentPath)
			// Send live metrics every 5 files to avoid UI overload
			if filesProcessed%5 == 0 {
				broker.SendMetric(filesProcessed, currentTokens, time.Since(startTime))
			}
		}

		onBlocked := func(path string) {
			broker.SendLog("WARN", "Filtered", path)
		}

		outputFile, tokens, filtered, err := core.GenerateProject(
			targetDir, outputDir, outName, req.CustomIgnores, req.Format, req.IncludePatterns,
			req.RemoveComments, req.RemoveEmptyLines, req.AiInstructions, req.CustomPrompt,
			req.AllowSecrets, onProgress, onBlocked,
		)

		// Send final metric
		broker.SendMetric(totalFilesProcessed, tokens, time.Since(startTime))

		if cleanupDir != "" {
			os.RemoveAll(cleanupDir)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		}

		if filtered > 0 {
			broker.SendLog("WARN", "Security", fmt.Sprintf("%d sensitive files were redacted.", filtered))
		}

		fmt.Fprintf(w, "Success! Generated: %s (Approx. %d Tokens)", outputFile, tokens)
	})

	// Route: Server-Sent Events (SSE) for Real-time Logs
	http.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		messageChan := make(chan string, 10000)
		broker.AddClient(messageChan)

		defer func() {
			broker.RemoveClient(messageChan)
		}()

		notify := r.Context().Done()
		for {
			select {
			case <-notify:
				return
			case msg := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
		}
	})

	// Route: Restore Project
	http.HandleFunc("/api/restore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req RestoreRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		startTime := time.Now()
		totalFilesProcessed := 0
		totalTokensProcessed := 0

		// Progress Callback for Real-time SSE Updates
		onProgress := func(filesProcessed, currentTokens int, currentPath string) {
			totalFilesProcessed = filesProcessed
			totalTokensProcessed = currentTokens
			broker.SendLog("INFO", "Restoring", currentPath)
			if filesProcessed%5 == 0 {
				broker.SendMetric(filesProcessed, currentTokens, time.Since(startTime))
			}
		}

		if err := core.RestoreProject(req.MdFilePath, req.DestinationDir, onProgress); err != nil {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		}

		// Send final metric
		broker.SendMetric(totalFilesProcessed, totalTokensProcessed, time.Since(startTime))

		fmt.Fprintf(w, "Success! Project restored into: %s", req.DestinationDir)
	})

	// Route: Open Folder in OS
	http.HandleFunc("/api/open-folder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if err := utils.OpenFolderInOS(req.Path); err != nil {
			http.Error(w, fmt.Sprintf("Failed to open folder: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Route: Custom File Explorer (Zero-Dependency)
	http.HandleFunc("/api/explore", handleExplore)
	http.HandleFunc("/api/quickaccess", handleQuickAccess)
	http.HandleFunc("/api/mkdir", handleMkdir)

	// Route: Safe Shutdown
	http.HandleFunc("/api/shutdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Shutting down...")
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}()
	})

	// Route: Explicit Disconnect
	http.HandleFunc("/api/disconnect", func(w http.ResponseWriter, r *http.Request) {
		heartbeatMu.Lock()
		// Zero-Bug Policy: Wait 5 seconds before shutting down.
		// If this was just a page refresh, the immediate heartbeat on page load will cancel this timer!
		shutdownTimer = time.AfterFunc(5*time.Second, func() {
			fmt.Println("Browser tab closed. Shutting down server safely...")
			os.Exit(0)
		})
		heartbeatMu.Unlock()
	})

	// Route: Heartbeat (Ping)
	http.HandleFunc("/api/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		heartbeatMu.Lock()
		if shutdownTimer != nil {
			shutdownTimer.Stop() // Cancel shutdown if page was just refreshed
			shutdownTimer = nil
		}
		lastHeartbeat = time.Now()
		heartbeatMu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	// Heartbeat Monitor Goroutine
	heartbeatMu.Lock()
	lastHeartbeat = time.Now() // Initial grace period
	heartbeatMu.Unlock()

	go func() {
		time.Sleep(15 * time.Second)
		for {
			time.Sleep(5 * time.Second)
			heartbeatMu.Lock()
			elapsed := time.Since(lastHeartbeat)
			heartbeatMu.Unlock()

			// Enterprise Policy: 1 HOUR timeout to act as a Garbage Collector for zombie processes.
			// Browser throttling (e.g., Chrome background tabs) won't kill the server prematurely.
			if elapsed > 1*time.Hour {
				fmt.Println("UI heartbeat lost for 1 hour (Zombie Process). Shutting down server safely...")
				os.Exit(0)
			}
		}
	}()

	// Open browser automatically without flashing CMD windows
	go func() {
		url := "http://localhost:" + port
		if err := utils.OpenBrowser(url); err != nil {
			fmt.Printf("Please open your browser and navigate to %s\n", url)
		}
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Fatal Error starting server: %v\n", err)
	}
}
