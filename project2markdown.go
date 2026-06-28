// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// outputFileSuffix defines the default naming convention suffix for the generated context file.
const outputFileSuffix = "_context.md"

// MaxLineBuffer defines a strict 20MB limit per single raw line to capture ultra-large compressed assets.
const MaxLineBuffer = 1024 * 1024 * 20

// excludeDirs lists the directory names (in strict lowercase) that should be skipped during traversal.
var excludeDirs = map[string]bool{
	".git":             true,
	".github":          true,
	"node_modules":     true,
	"vendor":           true,
	"bin":              true,
	"obj":              true,
	"binaries":         true,
	"deriveddatacache": true,
	"intermediate":     true,
	"saved":            true,
	".vs":              true,
}

// allowedExtensions defines a comprehensive list of supported file types for safe processing.
var allowedExtensions = map[string]bool {
	// Web & Scripting
	".php":  true,
	".json": true,
	".js":   true,
	".ts":   true,
	".jsx":  true,
	".tsx":  true,
	".css":  true,
	".html": true,
	
	// Systems & App Development
	".go":    true,
	".cpp":   true,
	".hpp":   true,
	".c":     true,
	".h":     true,
	".cs":    true,
	".java":  true,
	".py":    true,
	".rb":    true,
	".rs":    true,
	".swift": true,
	".kt":    true,
	
	// Game Engines & Custom Layouts (Unreal Engine / 3D Specs)
	".uproject": true,
	".uplugin":  true,
	".t3d":      true,
	".copy":     true,
	".usf":      true,
	".ush":      true,
	
	// Data & Documentation
	".md":   true,
	".txt":  true,
	".yaml": true,
	".yml":  true,
	".xml":  true,
	".ini":  true,
	".conf": true,
}

func main() {
	// 1. Retrieve the current working directory safely.
	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Fatal Error: Unable to retrieve current working directory: %v\n", err)
		return
	}
	projectName := filepath.Base(projectDir)
	outputFile := projectName + outputFileSuffix

	// 2. Create or truncate the output markdown file.
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Fatal Error: Unable to create output file: %v\n", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Warning: Error closing output file explicitly: %v\n", err)
		}
	}()

	// 3. Generate the directory structure tree safely.
	fmt.Printf("Project2Markdown: Processing root folder [%s]\n", projectName)
	fmt.Println("Generating secure directory tree...")
	_, _ = f.WriteString(fmt.Sprintf("# %s Project Context & Source Tree\n\n", projectName))
	_, _ = f.WriteString("## 1. Directory Structure\n```text\n")
	_, _ = f.WriteString(projectName + "/\n")
	
	generateTree(projectDir, projectDir, "", outputFile, f)
	_, _ = f.WriteString("```\n\n")

	_, _ = f.WriteString("==================================================\n")
	_, _ = f.WriteString("PROJECT STRUCTURE & EXISTING SOURCE CODE\n")
	_, _ = f.WriteString("==================================================\n\n")

	// 4. Recursive file traversal with chunk-based streaming safety.
	fmt.Println("Streaming and formatting source code safely...")
	err = filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil 
		}

		// Skip forbidden directories completely (Case-Insensitive Check).
		if d.IsDir() {
			if excludeDirs[strings.ToLower(d.Name())] {
				return filepath.SkipDir
			}
			return nil
		}

		// Zero-Bug Safety: Skip symlinks and hardware junction points to prevent endless circular routing loops.
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 || (info.Mode()&os.ModeDevice != 0) {
			return nil 
		}

		// Prevent the program from self-dumping
		lowerName := strings.ToLower(d.Name())
		if d.Name() == outputFile || 
			lowerName == "project2markdown.exe" || 
			lowerName == "project2markdown.go" || 
			strings.HasSuffix(lowerName, strings.ToLower(outputFileSuffix)) {
			return nil
		}

		// Process whitelisted extensions.
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if allowedExtensions[ext] {
			relPath, err := filepath.Rel(projectDir, path)
			if err != nil {
				return nil
			}
			relPath = filepath.ToSlash(relPath)

			langMarkdownTag := strings.TrimPrefix(ext, ".")

			_, _ = f.WriteString(fmt.Sprintf("[FILE: %s/%s]\n", projectName, relPath))
			_, _ = f.WriteString(fmt.Sprintf("```%s\n", langMarkdownTag)) 

			// Zero-Bug Safety: Stream contents safely.
			err = streamFileContents(path, f)
			if err != nil {
				_, _ = f.WriteString(fmt.Sprintf("// Error streaming file content completely: %v\n", err))
				fmt.Printf("Warning: File processed with issues [%s]: %v\n", d.Name(), err)
			}

			_, _ = f.WriteString("\n```\n")
			_, _ = f.WriteString("==================================================\n")
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Traversal Error: %v\n", err)
		return
	}

	fmt.Printf("\nSUCCESS! Context generated safely: %s\n", outputFile)
}

// generateTree acts as a robust recursive directory-mapping printer using an io.Writer interface.
func generateTree(rootDir, currentDir, prefix, outputFile string, w io.Writer) {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return 
	}

	var validEntries []fs.DirEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		lowerName := strings.ToLower(entry.Name())
		
		if !excludeDirs[lowerName] && 
			info.Mode()&os.ModeSymlink == 0 &&
			entry.Name() != outputFile && 
			lowerName != "project2markdown.exe" && 
			lowerName != "project2markdown.go" && 
			!strings.HasSuffix(lowerName, strings.ToLower(outputFileSuffix)) {
			validEntries = append(validEntries, entry)
		}
	}

	for i, entry := range validEntries {
		isLast := i == len(validEntries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		_, _ = fmt.Fprintf(w, "%s%s%s\n", prefix, connector, entry.Name())

		if entry.IsDir() {
			indent := "│   "
			if isLast {
				indent = "    "
			}
			generateTree(rootDir, filepath.Join(currentDir, entry.Name()), prefix+indent, outputFile, w)
		}
	}
}

// streamFileContents opens the source file and safely streams its content line-by-line.
func streamFileContents(srcPath string, dstFile *os.File) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, MaxLineBuffer) // Supports up to 20MB per single line for safety.

	for scanner.Scan() {
		line := scanner.Text()
		
		if strings.Contains(line, "```") {
			line = strings.ReplaceAll(line, "```", "` ` `")
		}
		
		_, err = dstFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		// If a line exceeds our max buffer, write a clear structural warning to the MD asset.
		if err == bufio.ErrTooLong {
			_, _ = dstFile.WriteString(fmt.Sprintf("\n// [Zero-Bug Warning: Single raw line exceeded safe threshold of %d MB]\n", MaxLineBuffer/(1024*1024)))
		}
		return err
	}

	return nil
}