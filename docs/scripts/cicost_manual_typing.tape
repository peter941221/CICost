Output docs/assets/cicost-cli-demo-v4.gif

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

Type "echo 'Demo: cicost version'"
Enter
Sleep 300ms

Type ".\\cicost.exe version"
Enter
Sleep 900ms

Type "echo 'Demo: cicost help'"
Enter
Sleep 300ms

Type ".\\cicost.exe help"
Enter
Sleep 900ms

Type "echo 'Demo: policy lint'"
Enter
Sleep 300ms

Type ".\\cicost.exe policy lint --policy .cicost.policy.yml.example"
Enter
Sleep 1000ms

Type "echo 'Demo: policy explain'"
Enter
Sleep 300ms

Type ".\\cicost.exe policy explain"
Enter
Sleep 900ms

Type "echo 'Demo: merged config'"
Enter
Sleep 300ms

Type ".\\cicost.exe config show"
Enter
Sleep 1200ms
