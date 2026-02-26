package cmd

import "fmt"

var (
	version = "0.1.0-dev"
	commit  = "none"
	builtAt = "unknown"
)

func runVersion(_ []string) error {
	fmt.Printf("cicost %s (commit: %s, built: %s)\n", version, commit, builtAt)
	return nil
}
