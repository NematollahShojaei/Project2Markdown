<div align="center">

<img width="120" height="120" alt="P2M Logo" src="https://github.com/user-attachments/assets/84b3c9a4-79ab-4a0a-9d54-0c5b72828e3d" />

# ⚡ Project2Markdown

### The two-way bridge between your codebase and AI.

**Compress an entire repository into one AI-ready file. Let the AI edit it. Restore it back into real code.**

<br/>

[![License: MIT](https://img.shields.io/badge/License-MIT-000000?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Zero Dependency](https://img.shields.io/badge/Dependencies-Zero-16A34A?style=for-the-badge)](#)
[![Version](https://img.shields.io/badge/Version-2.0.1-FF007F?style=for-the-badge)](#)

<br/>

[**Quick Start**](#quick-start) · [**GUI Mode**](#gui-mode) · [**CLI Mode**](#cli-mode) · [**Restore Mode**](#restore-mode) · [**Features**](#why-p2m)

</div>

<br/>

---

## 🤔 The Problem

You want an LLM to understand your entire codebase. Today that means manually copy-pasting files, fighting token limits, and losing structure. And once the AI hands back new code — you're pasting it back in, file by file.

**P2M turns that into one command in, one command out.**

```
project/         ──▶  p2m --cli  ──▶  context.md  ──▶  🤖 paste into Claude / GPT / Gemini
   ├─ src/                                                        │
   ├─ config/                                                     ▼
   └─ ...        ◀──  p2m --restore  ◀──  ai_output.xml  ◀──  AI rewrites your project
```

<br/>

<a id="quick-start"></a>
## 🚀 Quick Start

No cloning, no building — run it instantly from anywhere.

**🍎 macOS / Linux**
```bash
curl -fsSL https://raw.githubusercontent.com/nematollahshojaei/project2markdown/main/install.sh | bash
```

**🪟 Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/nematollahshojaei/project2markdown/main/install.ps1 | iex
```

**🐹 Go installed (cross-platform)**
```bash
go run github.com/nematollahshojaei/project2markdown/v2/cmd/p2m@latest --cli
```

> Once installed, type **`p2m`** to launch the Web UI, or **`p2m --cli`** to generate context directly in your terminal. That's it.

<br/>

<a id="gui-mode"></a>
## 🖥️ GUI Mode
<img align="right" width="360" alt="GUI preview" src="https://github.com/user-attachments/assets/b1ae9e95-c971-4e29-89c2-421d79cd6887" />

Run `p2m` with no flags and a local **Glassmorphism dashboard** opens in your browser automatically. Built for people who'd rather click than type.

- 📊 **Live metrics** — tokens/sec, elapsed time, file count, all updating in real time
- 🗂️ **Smart file explorer** — browse local drives or drop in a remote GitHub URL
- 🎛️ **One-click toggles** — strip comments, strip empty lines, custom AI prompts, and **Allow Sensitive Files**
- 📂 **Smart Output & Open Folder** — Auto-routes read-only paths to `Documents/Project2Markdown` and highlights the generated file in your OS natively.

<br clear="right"/>

<a id="cli-mode"></a>
## 💻 CLI Mode

For terminal-first developers. Every example below is copy-paste ready.

> 💡 **Tip:** Run `p2m --help` at any time to see the full list of available flags, options, and examples directly in your terminal.

<details open>
<summary><b>Basic generation (Markdown output)</b></summary>

```bash
p2m --cli
```
</details>

<details>
<summary><b>Advanced minification — XML, comments & empty lines stripped</b></summary>

```bash
p2m --cli --format=xml --remove-comments --remove-empty-lines
```
</details>

<details>
<summary><b>Process a remote GitHub repo — no <code>git clone</code> needed</b></summary>

```bash
p2m --cli --remote=yamadashy/repomix --format=json
```
</details>

<details>
<summary><b>Include only specific files or folders</b></summary>

```bash
p2m --cli --include="*.go, src/, config.json"
```
</details>

<details>
<summary><b>Specify output directory & allow sensitive files (.env, keys)</b></summary>

```bash
p2m --cli --output="./my_exports" --allow-secrets
```
</details>

<br/>

<a id="restore-mode"></a>
## 🔄 Restore Mode

Got a full project back from an AI as one Markdown/XML/JSON file? Rebuild every file and folder in one shot.

```bash
p2m --restore=my_project_context.xml
```

🛡️ **Bulletproof Path Traversal Protection** — Malicious or hallucinated paths (`../../etc/passwd`, `C:\Windows`, etc.) are mathematically blocked (LFI/RFI safe) before a single file is written, even on Root directories.

<br/>

<a id="why-p2m"></a>
## ✨ Why P2M

| | |
|---|---|
| ⚡ **Concurrency Engine** | Uses every CPU core via a worker pool — thousands of files processed in milliseconds |
| 🧠 **Built-in Tokenizer** | Zero-dependency heuristic tokenizer — accurate AI token counts, no heavy ML libs |
| 🌐 **Remote Fetching** | Point it at any public GitHub URL — no local clone required |
| 🤖 **AI-Optimized Formats** | Export as **Markdown**, **XML** *(best for Claude)*, or **JSON** *(best for APIs)* |
| 🛡️ **Universal Ignore Engine** | Auto-respects `.gitignore`, `.p4ignore`, and `.p2mignore`. Safely redacts `.env` and private keys by default. |
| 🔒 **Enterprise Security** | Hardened against SSRF, ZipSlip, LFI/RFI, and Memory/Goroutine Leaks. |
| 📂 **Smart Output Routing** | Automatically redirects outputs from Read-Only directories to a safe `Documents/Project2Markdown` workspace. |
| 🎮 **UE5-Native Parsing** | Understands `.uproject`, `.uplugin`, `.t3d`, `.usf`, `.ush` — skips `Saved`, `Intermediate`, `Binaries`, `DerivedDataCache` |

<br/>

## 🛠️ Building from Source

<details>
<summary>Manual build instructions (for contributors)</summary>

<br/>

**1. Clone the repository**
```bash
git clone https://github.com/nematollahshojaei/project2markdown.git
cd project2markdown
```

**2. Build the CLI edition**
```bash
go build -ldflags="-s -w" -o p2m-cli ./cmd/p2m
```

**3. Build the GUI edition** *(Windows, hides the console window)*
```bash
go build -ldflags="-s -w -H=windowsgui" -o p2m-gui.exe ./cmd/p2m
```

</details>

<br/>

---

<div align="center">

### 📄 License

Released under the **MIT License** — free for personal and commercial use.

<br/>

**Developed with 💙 by [Nematollah Shojaei](https://www.r3dhills.com)**
*Unreal Engine Technical Specialist & Software Engineer*

[🌐 Website](https://www.r3dhills.com) · [✉️ Email](mailto:nemat@r3dhills.com) · [💼 LinkedIn](https://linkedin.com/in/NematollahShojaei)

<br/>

⭐ **If P2M saves you time, consider starring the repo — it genuinely helps.**

</div>
