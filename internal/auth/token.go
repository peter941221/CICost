package auth

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/peter941221/CICost/internal/config"
)

var ErrNoToken = errors.New("no github token found")

// ResolveToken returns token in priority:
// explicit -> GITHUB_TOKEN -> GH_TOKEN -> gh auth token -> config token.
func ResolveToken(explicit string, cfgToken string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		return v, nil
	}
	if v := os.Getenv("GH_TOKEN"); v != "" {
		return v, nil
	}
	if v := tokenFromGHCLI(); v != "" {
		return v, nil
	}
	if v := strings.TrimSpace(config.Expand(cfgToken)); v != "" {
		return v, nil
	}
	return "", ErrNoToken
}

func tokenFromGHCLI() string {
	cmd := exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
