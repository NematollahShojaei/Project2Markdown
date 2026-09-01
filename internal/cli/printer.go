// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package cli

import (
	"fmt"
)

// ANSI Color Codes for Enterprise CLI styling
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
)

// PrintSummaryBox draws a beautiful, structured summary box in the terminal.
func PrintSummaryBox(title, path string, tokens int, filtered int) {
	fmt.Println()
	fmt.Printf("%s╭────────────────────────────────────────────────────────────╮%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s│ %s%-58s%s │%s\n", ColorGreen, ColorGreen, title, ColorGreen, ColorReset)
	fmt.Printf("%s├────────────────────────────────────────────────────────────┤%s\n", ColorGreen, ColorReset)

	pathStr := path
	if len(pathStr) > 48 {
		pathStr = "..." + pathStr[len(pathStr)-45:]
	}
	fmt.Printf("%s│ %sPath:   %-50s%s │%s\n", ColorGreen, ColorReset, pathStr, ColorGreen, ColorReset)

	if tokens > 0 {
		fmt.Printf("%s│ %sTokens: %-50d%s │%s\n", ColorGreen, ColorReset, tokens, ColorGreen, ColorReset)
	}
	if filtered > 0 {
		fmt.Printf("%s│ %sFiltered: %-48d%s │%s\n", ColorYellow, ColorReset, filtered, ColorYellow, ColorReset)
	}

	fmt.Printf("%s╰────────────────────────────────────────────────────────────╯%s\n", ColorGreen, ColorReset)
}

// CustomHelp overrides the default flag.Usage to provide an Enterprise-grade, colorized help menu.
func CustomHelp() {
	banner := `
 ____  ____  __  __ 
(  _ \(___ \(  \/  )
 )___/ / __/ )    ( 
(__)  (____)(_/\/\_)
`
	fmt.Println(ColorCyan + banner + ColorReset)
	fmt.Println(ColorGreen + "Project2Markdown (P2M) - Enterprise AI Context Engine" + ColorReset)
	fmt.Println(ColorGray + "Version: 2.0.1 | Zero-Dependency" + ColorReset)

	fmt.Println("\n" + ColorYellow + "USAGE:" + ColorReset)
	fmt.Println("  p2m [flags]")

	fmt.Println("\n" + ColorYellow + "CORE FLAGS:" + ColorReset)
	fmt.Printf("  %s--cli%s                  Run in CLI mode (bypasses the Web UI)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--restore <path>%s       Restore a project from a generated context file\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--remote <url>%s         Process a remote GitHub repository (e.g., owner/repo)\n", ColorCyan, ColorReset)

	fmt.Println("\n" + ColorYellow + "GENERATION OPTIONS:" + ColorReset)
	fmt.Printf("  %s--format <type>%s        Output format: md, xml, json (Default: md)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--output <path>%s        Specify the output directory (Default: current directory)\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--include <rules>%s      Comma-separated list of files/folders to strictly include\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--remove-comments%s      Strip comments from code to save AI tokens\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--remove-empty-lines%s   Strip empty lines to save AI tokens\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--ai-instructions%s      Inject strict restore instructions for AI at the top\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--custom-prompt%s        Inject a custom prompt at the top of the file\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--allow-secrets%s        Allow inclusion of sensitive files (.env, keys)\n", ColorCyan, ColorReset)

	fmt.Println("\n" + ColorYellow + "EXAMPLES:" + ColorReset)
	fmt.Println("  p2m --cli")
	fmt.Println("  p2m --cli --format=xml --remove-comments")
	fmt.Println("  p2m --cli --remote=NematollahShojaei/Project2Markdown")
	fmt.Println("  p2m --restore=my_project_context.xml")
	fmt.Println()
}
