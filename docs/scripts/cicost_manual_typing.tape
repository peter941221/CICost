Output docs/assets/cicost-cli-demo.gif

Set Shell "pwsh"
Set FontSize 18
Set Width 1100
Set Height 700
Set Framerate 30
Set Theme "Builtin Solarized Dark"
Set Padding 20

Type "cd 'C:\\AI Projects\\Fun Stuff\\IndieHacker\\CICost'"
Enter
Sleep 300ms

Type "Write-Host '[1/5] CLI version' -ForegroundColor Yellow"
Enter
Sleep 300ms

Type ".\\cicost.exe version"
Enter
Sleep 900ms

Type "Write-Host '[2/5] Command overview (help)' -ForegroundColor Yellow"
Enter
Sleep 300ms

Type ".\\cicost.exe help | Select-Object -First 14"
Enter
Sleep 900ms

Type "Write-Host '[3/5] Effective merged config (config show)' -ForegroundColor Yellow"
Enter
Sleep 300ms

Type ".\\cicost.exe config show | Select-Object -First 20"
Enter
Sleep 1000ms

Type "Write-Host '[4/5] Policy file lint' -ForegroundColor Yellow"
Enter
Sleep 300ms

Type ".\\cicost.exe policy lint --policy .cicost.policy.yml.example"
Enter
Sleep 900ms

Type "Write-Host '[5/5] Policy expression reference' -ForegroundColor Yellow"
Enter
Sleep 300ms

Type ".\\cicost.exe policy explain | Select-Object -First 16"
Enter
Sleep 1200ms
