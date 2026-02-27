param(
  [string]$InputGif = "docs/assets/cicost-cli-demo-v6.gif",
  [string]$SubtitleFile = "docs/scripts/cicost_demo_v7.srt",
  [string]$OutputGif = "docs/assets/cicost-cli-demo-v7.gif"
)

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot "..\.."))

$tempMp4 = "docs/assets/_cicost-cli-demo-v7-sub.mp4"
$palette = "docs/assets/_cicost-cli-demo-v7-palette.png"

$style = "FontName=Arial,FontSize=20,PrimaryColour=&H00FFFFFF,OutlineColour=&H00000000,BackColour=&H80000000,BorderStyle=3,MarginV=22"
$subFilter = "subtitles='$SubtitleFile':force_style='$style'"

ffmpeg -hide_banner -y -i $InputGif -vf $subFilter $tempMp4 | Out-Null
ffmpeg -hide_banner -y -i $tempMp4 -vf "fps=15,scale=1100:-1:flags=lanczos,palettegen" $palette | Out-Null
ffmpeg -hide_banner -y -i $tempMp4 -i $palette -filter_complex "fps=15,scale=1100:-1:flags=lanczos[x];[x][1:v]paletteuse" $OutputGif | Out-Null

Remove-Item $tempMp4, $palette -Force -ErrorAction SilentlyContinue

$sizeMb = [Math]::Round(((Get-Item $OutputGif).Length / 1MB), 2)
Write-Host "Generated: $OutputGif ($sizeMb MB)" -ForegroundColor Green
