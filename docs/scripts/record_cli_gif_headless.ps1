param(
  [string]$Tape = "docs/scripts/cicost_manual_typing.tape"
)

$ErrorActionPreference = "Stop"

Set-Location (Resolve-Path (Join-Path $PSScriptRoot "..\.."))

$vhs = Join-Path $env:USERPROFILE "go\bin\vhs.exe"
if (-not (Test-Path $vhs)) {
  Write-Host "Installing vhs..." -ForegroundColor Yellow
  go install github.com/charmbracelet/vhs@latest
}

if (-not (Get-Command ttyd -ErrorAction SilentlyContinue)) {
  Write-Host "Installing ttyd..." -ForegroundColor Yellow
  scoop install ttyd
}

Write-Host "Building cicost.exe..." -ForegroundColor Yellow
go build -o cicost.exe .

Write-Host "Rendering GIF from tape (headless)..." -ForegroundColor Yellow
& $vhs $Tape

Write-Host "Done." -ForegroundColor Green
