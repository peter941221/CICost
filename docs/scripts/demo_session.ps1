param(
  [Parameter(Mandatory = $true)]
  [string]$RepoRoot
)

Set-Location $RepoRoot
$Host.UI.RawUI.WindowTitle = "CICost CLI Demo"

Write-Host "CICost CLI Demo Session" -ForegroundColor Yellow
Write-Host "Repository: $RepoRoot" -ForegroundColor DarkGray
Start-Sleep -Milliseconds 800

if (-not (Test-Path ".\cicost.exe")) {
  Write-Host "`n$ go build -o cicost.exe ." -ForegroundColor Cyan
  go build -o cicost.exe .
  Start-Sleep -Milliseconds 800
}

Write-Host "`n$ .\cicost.exe version" -ForegroundColor Cyan
.\cicost.exe version
Start-Sleep -Milliseconds 1100

Write-Host "`n$ .\cicost.exe help" -ForegroundColor Cyan
.\cicost.exe help
Start-Sleep -Milliseconds 1300

Write-Host "`n$ .\cicost.exe policy explain" -ForegroundColor Cyan
.\cicost.exe policy explain
Start-Sleep -Milliseconds 1300

Write-Host "`n$ .\cicost.exe report --repo owner/repo --days 7 --format table" -ForegroundColor Cyan
.\cicost.exe report --repo owner/repo --days 7 --format table
Start-Sleep -Milliseconds 1800

exit
