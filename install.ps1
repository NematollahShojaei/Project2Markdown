# Copyright 2026 R3D HILLS. All Rights Reserved.
# Enterprise Zero-Install Script for Project2Markdown (Windows)

$ErrorActionPreference = "Stop"

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "  Installing Project2Markdown (P2M) Enterprise... " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

$BinaryName = "p2m-windows-cli.exe"
$DownloadUrl = "https://github.com/nematollahshojaei/project2markdown/releases/latest/download/$BinaryName"

# Install in the user's profile directory to avoid requiring Administrator privileges
$InstallDir = Join-Path $env:USERPROFILE ".p2m\bin"
$DestFile = Join-Path $InstallDir "p2m.exe"

# 1. Create directory if it doesn't exist
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

Write-Host "Downloading latest release from GitHub..." -ForegroundColor White

# 2. Download the binary
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $DestFile -UseBasicParsing
} catch {
    Write-Host "Error: Failed to download the binary. Please check your internet connection or ensure a release exists on GitHub." -ForegroundColor Red
    exit 1
}

# 3. Add to User PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notmatch [regex]::Escape($InstallDir)) {
    Write-Host "Adding $InstallDir to your Environment PATH..." -ForegroundColor Yellow
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
    
    # Update current session PATH so it works immediately without restarting the terminal
    $env:PATH = "$env:PATH;$InstallDir" 
}

Write-Host "==================================================" -ForegroundColor Green
Write-Host "  SUCCESS! P2M has been installed successfully.   " -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
Write-Host "You can now run it from anywhere using the command:" -ForegroundColor White
Write-Host "  p2m --cli" -ForegroundColor Cyan
Write-Host ""
Write-Host "(Note: If the command is not recognized in your current VS Code terminal, please open a new terminal window)" -ForegroundColor Yellow