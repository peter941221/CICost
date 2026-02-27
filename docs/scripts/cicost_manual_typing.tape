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
Sleep 700ms

Type "go build -o cicost.exe ."
Enter
Sleep 1500ms

Type ".\\cicost.exe version"
Enter
Sleep 1200ms

Type ".\\cicost.exe help"
Enter
Sleep 1400ms

Type ".\\cicost.exe policy explain"
Enter
Sleep 1500ms

Type ".\\cicost.exe config show"
Enter
Sleep 1500ms

Type ".\\cicost.exe policy lint --policy .cicost.policy.yml.example"
Enter
Sleep 1800ms
