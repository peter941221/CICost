param(
  [int]$DurationSec = 18,
  [int]$Width = 1100,
  [string]$OutputGif = "docs/assets/cicost-cli-demo.gif"
)

$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
Set-Location $repoRoot

New-Item -ItemType Directory -Path "docs/assets" -Force | Out-Null

Write-Host "Building cicost.exe..." -ForegroundColor Yellow
go build -o cicost.exe .

$sessionScript = Join-Path $PSScriptRoot "demo_session.ps1"
$sessionProc = Start-Process powershell -ArgumentList @(
  "-NoLogo",
  "-NoProfile",
  "-ExecutionPolicy", "Bypass",
  "-File", "`"$sessionScript`"",
  "-RepoRoot", "`"$repoRoot`""
) -PassThru

Start-Sleep -Seconds 1

$mp4 = Join-Path $repoRoot "docs/assets/_cicost-cli-demo.mp4"
$palette = Join-Path $repoRoot "docs/assets/_palette.png"
$gifPath = Join-Path $repoRoot $OutputGif

Write-Host "Recording desktop for $DurationSec seconds..." -ForegroundColor Yellow
& ffmpeg -hide_banner -y -f gdigrab -framerate 10 -t $DurationSec -i desktop -vf "scale=${Width}:-1:flags=lanczos" $mp4 | Out-Null

Write-Host "Encoding GIF..." -ForegroundColor Yellow
& ffmpeg -hide_banner -y -i $mp4 -vf "fps=8,scale=${Width}:-1:flags=lanczos,palettegen" $palette | Out-Null
& ffmpeg -hide_banner -y -i $mp4 -i $palette -filter_complex "fps=8,scale=${Width}:-1:flags=lanczos[x];[x][1:v]paletteuse" $gifPath | Out-Null

if (-not $sessionProc.HasExited) {
  $null = $sessionProc | Wait-Process -Timeout 5 -ErrorAction SilentlyContinue
}

Remove-Item $mp4, $palette -Force -ErrorAction SilentlyContinue

$sizeMb = [Math]::Round(((Get-Item $gifPath).Length / 1MB), 2)
Write-Host "GIF generated: $gifPath ($sizeMb MB)" -ForegroundColor Green
