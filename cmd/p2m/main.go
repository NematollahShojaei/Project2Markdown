// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nematollahshojaei/project2markdown/v2/internal/cli"
	"github.com/nematollahshojaei/project2markdown/v2/internal/core"
	"github.com/nematollahshojaei/project2markdown/v2/internal/remote"
	"github.com/nematollahshojaei/project2markdown/v2/internal/server"
	"github.com/nematollahshojaei/project2markdown/v2/web"
)

func main() {
	// Override default flag usage with our Enterprise Help Menu
	flag.Usage = cli.CustomHelp

	// Parse CLI flags
	restoreFlag := flag.String("restore", "", "Path to the markdown file to restore the project from")
	cliFlag := flag.Bool("cli", false, "Run in CLI mode to generate context in the current directory without UI")
	formatFlag := flag.String("format", "md", "Output format for CLI mode: md, xml, json")
	outputFlag := flag.String("output", "", "Specify the output directory (default: source directory)")
	includeFlag := flag.String("include", "", "Comma-separated list of files/folders to strictly include")
	removeCommentsFlag := flag.Bool("remove-comments", false, "Remove comment lines to save tokens")
	removeEmptyLinesFlag := flag.Bool("remove-empty-lines", false, "Remove empty lines to save tokens")
	aiInstructionsFlag := flag.Bool("ai-instructions", false, "Inject strict restore instructions for AI")
	allowSecretsFlag := flag.Bool("allow-secrets", false, "Allow inclusion of sensitive files (.env, keys)")
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

		// Zero-Bug Policy: Apply the output flag to the restore engine as well
		err := core.RestoreProject(*restoreFlag, *outputFlag, onProgress)
		spinner.Stop()

		if err != nil {
			fmt.Printf("\n%sFatal Error: %v%s\n", cli.ColorRed, err, cli.ColorReset)
		} else {
			cli.PrintSummaryBox("SUCCESS! Project Restored", *restoreFlag, 0, 0)
		}
		return
	}

	// Route 2: Fast CLI Generation Mode
	if *cliFlag || *remoteFlag != "" {
		targetDir := ""
		cleanupDir := ""
		outName := ""
		outputDir := *outputFlag

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

			// Zero-Bug Policy: Prevent saving the output file inside the temporary directory
			if outputDir == "" {
				cwd, _ := os.Getwd()
				outputDir = cwd
			}

			parts := strings.Split(strings.Trim(strings.TrimPrefix(strings.TrimPrefix(*remoteFlag, "https://github.com/"), "github.com/"), "/"), "/")
			repoName := "repo"
			if len(parts) >= 2 {
				repoName = parts[1]
			}
			outName = repoName + "_context." + *formatFlag
		}

		spinner := cli.NewSpinner()
		fmt.Printf("%sProject2Markdown: Processing folder in %s format%s\n", cli.ColorCyan, strings.ToUpper(*formatFlag), cli.ColorReset)

		// UX Standard: Warn the user explicitly if they bypass the security blacklist
		if *allowSecretsFlag {
			fmt.Printf("%s⚠️ WARNING: Sensitive files inclusion is enabled. Ensure you do not leak production secrets to AI.%s\n", cli.ColorYellow, cli.ColorReset)
		}

		spinner.Start("Generating Context")

		onProgress := func(filesProcessed, currentTokens int, currentPath string) {
			spinner.Update(filesProcessed, currentTokens)
		}

		onBlocked := func(path string) {
			fmt.Printf("\r\033[K%s[FILTERED]%s %s (Sensitive File)\n", cli.ColorYellow, cli.ColorReset, path)
		}

		outputFile, tokens, filtered, err := core.GenerateProject(
			targetDir, outputDir, outName, "", *formatFlag, *includeFlag,
			*removeCommentsFlag, *removeEmptyLinesFlag, *aiInstructionsFlag,
			*customPromptFlag, *allowSecretsFlag, onProgress, onBlocked,
		)

		spinner.Stop()

		if cleanupDir != "" {
			os.RemoveAll(cleanupDir)
		}

		if err != nil {
			fmt.Printf("\n%sFatal Error: %v%s\n", cli.ColorRed, err, cli.ColorReset)
		} else {
			cli.PrintSummaryBox("SUCCESS! Context Generated", outputFile, tokens, filtered)
		}
		return
	}

	// Route 3: Default to Web UI (Double-click friendly)
	// If no flags are provided, or --ui is explicitly called, launch the UI.
	_ = uiFlag // Ignored, as UI is now the default fallback

	server.StartUI("8989", web.UIHTML)
}
