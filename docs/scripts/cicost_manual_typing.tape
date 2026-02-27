Output docs/assets/cicost-cli-demo-v6.gif

Set Shell "cmd"
Set FontSize 18
Set Width 1100
Set Height 700
Set Framerate 30
Set Theme "Builtin Solarized Dark"
Set Padding 20

Type 'cd /d "C:\\AI Projects\\Fun Stuff\\IndieHacker\\CICost"'
Enter
Sleep 300ms

Type "echo Goal: identify CI cost hotspots in 30 seconds"
Enter
Sleep 600ms

Type "echo Scan latest workflow data"
Enter
Sleep 300ms

Type ".\\cicost.exe scan --repo cli/cli --days 3 --workers 4"
Enter
Sleep 1200ms

Type "echo Cost summary"
Enter
Sleep 300ms

Type ".\\cicost.exe report --repo cli/cli --days 3 --format table"
Enter
Sleep 1000ms

Type "echo Top expensive workflows"
Enter
Sleep 300ms

Type ".\\cicost.exe hotspots --repo cli/cli --days 3 --group-by workflow --top 3 --sort cost --format table"
Enter
Sleep 900ms

Type "echo Policy gate result"
Enter
Sleep 300ms

Type ".\\cicost.exe policy check --repo cli/cli --days 3 --policy .cicost.policy.yml.example"
Enter
Sleep 800ms

Type "echo Next step: tune policy thresholds for your repo/team"
Enter
Sleep 1200ms
