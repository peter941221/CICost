package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peter941221/CICost/internal/config"
)

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	yesFlag := fs.Bool("yes", false, "使用默认值，不交互")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg := config.Default()
	if *yesFlag {
		path, err := config.SaveUserConfig(cfg)
		if err != nil {
			return err
		}
		fmt.Printf("Config written: %s\n", path)
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("CICost init")
	fmt.Println("Press Enter to accept default values.")

	cfg.Repos = promptList(reader, "Default repositories (owner/repo, comma-separated)", "")
	if len(cfg.Repos) == 0 {
		cfg.Repos = []string{}
	}
	cfg.Scan.Days = promptInt(reader, "Default scan days", cfg.Scan.Days)
	cfg.Scan.Workers = promptInt(reader, "Default worker count", cfg.Scan.Workers)
	cfg.Budget.Monthly = promptFloat(reader, "Monthly budget threshold (USD)", cfg.Budget.Monthly)

	path, err := config.SaveUserConfig(cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Config written: %s\n", path)
	return nil
}

func promptInt(r *bufio.Reader, title string, def int) int {
	fmt.Printf("%s [%d]: ", title, def)
	in, _ := r.ReadString('\n')
	in = strings.TrimSpace(in)
	if in == "" {
		return def
	}
	var n int
	if _, err := fmt.Sscanf(in, "%d", &n); err != nil || n <= 0 {
		return def
	}
	return n
}

func promptFloat(r *bufio.Reader, title string, def float64) float64 {
	fmt.Printf("%s [%.2f]: ", title, def)
	in, _ := r.ReadString('\n')
	in = strings.TrimSpace(in)
	if in == "" {
		return def
	}
	var f float64
	if _, err := fmt.Sscanf(in, "%f", &f); err != nil || f <= 0 {
		return def
	}
	return f
}

func promptList(r *bufio.Reader, title string, def string) []string {
	if def == "" {
		fmt.Printf("%s: ", title)
	} else {
		fmt.Printf("%s [%s]: ", title, def)
	}
	in, _ := r.ReadString('\n')
	in = strings.TrimSpace(in)
	if in == "" {
		in = def
	}
	if in == "" {
		return nil
	}
	parts := strings.Split(in, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
