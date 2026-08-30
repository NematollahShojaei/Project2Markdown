// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nematollahshojaei/project2markdown/internal/cli"
	"github.com/nematollahshojaei/project2markdown/internal/core"
	"github.com/nematollahshojaei/project2markdown/internal/remote"
	"github.com/nematollahshojaei/project2markdown/internal/server"
	"github.com/nematollahshojaei/project2markdown/web"
)

func main() {
	// Override default flag usage with our Enterprise Help Menu
	flag.Usage = cli.CustomHelp

	// Parse CLI flags
	restoreFlag := flag.String("restore", "", "Path to the markdown file to restore the project from")
	cliFlag := flag.Bool("cli", false, "Run in CLI mode to generate context in the current directory without UI")
	formatFlag := flag.String("format", "md", "Output format for CLI mode: md, xml, json")
	includeFlag := flag.String("include", "", "Comma-separated list of files/folders to strictly include")
	removeCommentsFlag := flag.Bool("remove-comments", false, "Remove comment lines to save tokens")
	removeEmptyLinesFlag := flag.Bool("remove-empty-lines", false, "Remove empty lines to save tokens")
	aiInstructionsFlag := flag.Bool("ai-instructions", false, "Inject strict restore instructions for AI")
	customPromptFlag := flag.String("custom-prompt", "", "Inject a custom prompt at the top of the file")
	uiFlag := flag.Bool("ui", false, "Launch the visual web interface (Default behavior)")
	remoteFlag := flag.String("remote", "", "GitHub repository URL or owner/repo to process remotely")
	flag.Parse()

	// Route 1: Reverse Generation Engine (CLI)
	if *restoreFlag != "" {
		spinner := cli.NewSpinner()
		fmt.Printf("%sProject2Markdown: Initiating Reverse Generation from [%s]%s\n", cli.ColorCyan, *restoreFlag, cli.ColorReset)
		spinner.Start("Restoring Project")

		onProgress := func(filesProcessed, currentTokens int, currentPath string) {
			spinner.Update(filesProcessed, currentTokens)
		}

		err := core.RestoreProject(*restoreFlag, "", onProgress)
		spinner.Stop()

		if err != nil {
			fmt.Printf("\n%sFatal Error: %v%s\n", cli.ColorRed, err, cli.ColorReset)
		} else {
			cli.PrintSummaryBox("SUCCESS! Project Restored", *restoreFlag, 0)
		}
		return
	}

	// Route 2: Fast CLI Generation Mode
	if *cliFlag || *remoteFlag != "" {
		targetDir := ""
		cleanupDir := ""
		outName := ""

		if *remoteFlag != "" {
			fmt.Printf("Downloading remote repository: %s...\n", *remoteFlag)

			onDownloadProgress := func(downloaded, total int64) {
				mb := float64(downloaded) / (1024 * 1024)
				if total > 0 {
					totalMB := float64(total) / (1024 * 1024)
					fmt.Printf("\r\033[KDownloading... %.1f MB / %.1f MB", mb, totalMB)
				} else {
					fmt.Printf("\r\033[KDownloading... %.1f MB", mb)
				}
			}

			extractedPath, tmpDir, err := remote.DownloadAndExtractRepo(*remoteFlag, onDownloadProgress)
			fmt.Println() // Move to the next line after download finishes

			if err != nil {
				fmt.Printf("Fatal Error downloading repo: %v\n", err)
				return
			}
			targetDir = extractedPath
			cleanupDir = tmpDir

			cwd, _ := os.Getwd()
			parts := strings.Split(strings.Trim(strings.TrimPrefix(strings.TrimPrefix(*remoteFlag, "https://github.com/"), "github.com/"), "/"), "/")
			repoName := "repo"
			if len(parts) >= 2 {
				repoName = parts[1]
			}
			outName = filepath.Join(cwd, repoName+"_context."+*formatFlag)
		}

		spinner := cli.NewSpinner()
		fmt.Printf("%sProject2Markdown: Processing folder in %s format%s\n", cli.ColorCyan, strings.ToUpper(*formatFlag), cli.ColorReset)
		spinner.Start("Generating Context")

		onProgress := func(filesProcessed, currentTokens int, currentPath string) {
			spinner.Update(filesProcessed, currentTokens)
		}

		outputFile, tokens, err := core.GenerateProject(
			targetDir, outName, "", *formatFlag, *includeFlag,
			*removeCommentsFlag, *removeEmptyLinesFlag, *aiInstructionsFlag,
			*customPromptFlag, onProgress,
		)

		spinner.Stop()

		if cleanupDir != "" {
			os.RemoveAll(cleanupDir)
		}

		if err != nil {
			fmt.Printf("\n%sFatal Error: %v%s\n", cli.ColorRed, err, cli.ColorReset)
		} else {
			cli.PrintSummaryBox("SUCCESS! Context Generated", outputFile, tokens)
		}
		return
	}

	// Route 3: Default to Web UI (Double-click friendly)
	// If no flags are provided, or --ui is explicitly called, launch the UI.
	_ = uiFlag // Ignored, as UI is now the default fallback

	server.StartUI("8989", web.UIHTML)
}
