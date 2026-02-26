package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/peter941221/CICost/internal/config"
)

func runConfig(args []string) error {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	sub := "show"
	if fs.NArg() > 0 {
		sub = strings.ToLower(fs.Arg(0))
	}

	switch sub {
	case "show":
		cfg, err := config.LoadMerged(".cicost.yml")
		if err != nil {
			return err
		}
		b, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		fmt.Print(string(b))
		return nil
	case "edit":
		p, err := config.UserConfigPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				if _, err := config.SaveUserConfig(config.Default()); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		cmd := exec.Command("notepad", p)
		return cmd.Start()
	default:
		return fmt.Errorf("unknown subcommand %q, expected show|edit", sub)
	}
}
