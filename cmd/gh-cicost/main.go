package main

import (
	"fmt"
	"os"

	"github.com/peter941221/CICost/cmd"
)

func main() {
	if err := cmd.Execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(cmd.ExitCode(err))
	}
}
