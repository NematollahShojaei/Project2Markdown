# Project2Markdown (P2M) 🚀

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org/)

**Project2Markdown** is a high-performance, ultra-lightweight CLI utility written in Go[cite: 1]. It is engineered from the ground up to compress your entire project architecture and source code into a single, beautifully organized Markdown file optimized specifically for AI/LLM consumption (ChatGPT, Claude, Gemini).

Stop fighting token limits, context fragmentation, and out-of-memory crashes on massive repositories. 

---

## 📖 Official Guides & Tutorials

To get the absolute most out of this tool, check out our deep-dive engineering articles on **R3D HILLS**:

*   **[Core Tool Architecture]** [Discover how P2M handles gigabyte-scale repositories efficiently](https://r3dhills.com/project2markdown-ai-context-generator/)
*   **[Unreal Engine 5 Tutorial]** [Step-by-Step Guide: How to give your entire UE5 project to AI](https://r3dhills.com/how-to-give-ue5-project-to-ai/)

---

## ⚡ Key Features

| Feature | Description |
| :--- | :--- |
| **Stream-Based Architecture** | Uses a `bufio.Scanner` pipeline to read/write line-by-line[cite: 1]. Near-zero memory footprint, even on enterprise codebases. |
| **Directory Tree Generation** | Automatically injects an ASCII folder map at the top, giving the LLM an instant birds-eye view of your architecture[cite: 1]. |
| **Smart Exclusion Grid** | Out-of-the-box filtering for heavy directories (`.git`, `node_modules`, `Intermediate`, `Saved`)[cite: 1]. Avoids self-inclusion automatically. |
| **Markdown Safety Grid** | Automatically detects and sanitizes pre-existing triple backticks (```) inside code files to prevent formatting breaks[cite: 1]. |

---

## 🎮 Tailored for Unreal Engine 5 (UE5)

While P2M works flawlessly with any programming language, it features native optimizations designed specifically for game developers[cite: 1]. 

Visual scripting languages (like UE Blueprints and Maps) are traditionally stored as heavy binary assets (`.uasset`/`.umap`), making them unreadable for LLMs. P2M bridges this gap perfectly by natively parsing text-based metadata, shiders, and configuration files[cite: 1]:
*   **Supported Extensions:** `.uproject`, `.uplugin`, `.t3d`, `.usf`, `.ush`[cite: 1]
*   **Automatic Noise Reduction:** Automatically skips massive compiled engine folders like `Saved`, `Intermediate`, `Binaries`, and `DerivedDataCache`[cite: 1].

> **💡 The Blueprint Pipeline:** By combining Unreal Engine's native **Bulk Export** feature (exporting logic to `.T3D` or `.COPY` formats) with Project2Markdown, you can successfully pack your entire visual scripting architecture into a lightweight text context that Claude or ChatGPT can fully understand, refactor, and debug.

---

## 🛠️ Quick Start

### 1. Download the Binary
Head over to the **[Releases](https://github.com/nematollahshojaei/project2markdown/releases)** section and download the latest compiled executable for your operating system.

### 2. Execution
Drop the `project2markdown` executable directly into the root folder of the project you want to process.

*   **Windows:** Right-click and select **Run as Administrator**.
*   **Linux/macOS:** Run via terminal:
    ```bash
    chmod +x project2markdown
    ./project2markdown
    ```

### 3. Output
The tool will instantly generate a consolidated context file named `[Your_Project_Name]_context.md` in the same directory, completely optimized and ready to be dragged-and-dropped into any AI interface.

---

## 📄 License

This project is open-source and software released under the permissive **MIT License**. Feel free to modify, distribute, or integrate it into your automated CI/CD pipelines.

---

## 🤝 Connect & Network

Developed with 💙 by **Nematollah Shojaei** — Unreal Engine Technical Specialist & Software Engineer.

*   **Website:** [r3dhills.com](https://www.r3dhills.com)
*   **LinkedIn:** [linkedin.com/in/NematollahShojaei](https://linkedin.com/in/NematollahShojaei)
*   **Contact:** nemat@r3dhills.com