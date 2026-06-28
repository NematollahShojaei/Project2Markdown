# Project2Markdown 🚀

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://github.com/NematollahShojaei/project2markdown/pulls)
[![Maintained by R3D HILLS](https://img.shields.io/badge/Maintained%20by-R3D%20HILLS-black?style=flat-square)](https://r3dhills.com)

**Project2Markdown** is a high-performance, stream-based CLI tool written in Go that securely traverses your project directory, maps its structural hierarchy, and consolidates all relevant source code into a single, beautifully formatted Markdown file.

Built explicitly to streamline **LLM Context Feeding**, this tool allows you to pack your entire codebase into a clean, structured text asset that AI models (like Gemini, Claude, and ChatGPT) can instantly read, understand, and analyze without losing architectural context.

Developed and maintained by **R3D HILLS**.

---

## ✨ Key Features

- 📂 **Automated Directory Mapping:** Generates a visual, beautifully structured ASCII directory tree at the top of your context file.
- 🏎️ **Memory-Efficient Streaming:** Uses line-by-line chunked streaming via `bufio.Scanner` rather than buffering entire files into RAM—perfect for massive enterprise codebases.
- 🎮 **Unreal Engine & 3D Optimized:** First-class support for game developers and technical artists. It natively processes UE specs (`.uproject`, `.uplugin`, `.t3d`, `.usf`, `.ush`) while automatically ignoring massive engine-generated bloat folders.
- 🛡️ **Zero-Bug Safety Grid:**
  - **Circular Loop Protection:** Automatically skips symlinks and hardware junction points to prevent endless routing loops.
  - **Smart Self-Exclusion:** Prevents self-dumping by automatically ignoring its own source code, binary, and output files.
  - **Markdown Sanitization:** Safely escapes conflicting triple backticks (`` ``` ``) inside your source files to ensure the final output layout never breaks.
  - **Ultra-Large Asset Buffer:** Supports single raw lines up to **20MB** (ideal for heavily minified or compressed data arrays) with clear truncation warnings.

---

## 🛠️ Supported Environments & Extensions

Project2Markdown filters out binaries, images, and heavy build artifacts out of the box, focusing strictly on meaningful human-readable source code:

| Category | Supported Extensions |
| :--- | :--- |
| **Web & Scripting** | `.js`, `.ts`, `.jsx`, `.tsx`, `.php`, `.html`, `.css`, `.json` |
| **Systems & Mobile** | `.go`, `.cpp`, `.hpp`, `.c`, `.h`, `.cs`, `.java`, `.py`, `.rb`, `.rs`, `.swift`, `.kt` |
| **Game Engines (Unreal)** | `.uproject`, `.uplugin`, `.t3d`, `.copy`, `.usf`, `.ush` |
| **Data & Configuration** | `.md`, `.txt`, `.yaml`, `.yml`, `.xml`, `.ini`, `.conf` |

### 🚫 Automatically Ignored Directories
To keep your context lightweight, the following directories are completely skipped during traversal:
` .git ` • ` .github ` • ` node_modules ` • ` vendor ` • ` bin ` • ` obj ` • ` binaries ` • ` intermediate ` • ` saved ` • ` deriveddatacache ` • ` .vs `

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://go.dev/doc/install) (version 1.16 or higher recommended)

### Installation & Build

1. Clone the repository directly into your local machine:
   ```bash
   git clone https://github.com/NematollahShojaei/project2markdown.git
   cd project2markdown
   ```
2. Build the production-ready executable:
   ```bash
   go build -ldflags="-s -w" -o project2markdown.exe project2markdown.go
   ```

### **How to Use**

Simply place the generated executable (project2markdown.exe) or run the script directly inside the root directory of the project you want to process:

```bash
# Run via Go directly
go run project2markdown.go

# Or run the compiled binary
./project2markdown
```

### **Output**

The tool reads the folder name (e.g., MyAwesomeApp) and instantly drops a beautifully packed context file named:
```text
MyAwesomeApp_context.md
```

## **📂 Example Output Structure**

The generated context file follows a strict, highly scannable design:

```text
# MyAwesomeApp Project Context & Source Tree

## 1. Directory Structure
MyAwesomeApp/
├── internal/
│   └── parser.go
└── config.yaml

==================================================
PROJECT STRUCTURE & EXISTING SOURCE CODE
==================================================

[FILE: MyAwesomeApp/internal/parser.go]  
```
```go
package internal  
// Your streamed code contents seamlessly processed here... 
```

## **📝 License**

Distributed under the MIT License. See LICENSE for more information.  
Copyright 2026 **R3D HILLS**. All Rights Reserved.