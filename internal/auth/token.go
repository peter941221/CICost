package auth

import (
	"errors"
	"os"
)

var ErrNoToken = errors.New("no github token found")

// ResolveToken returns token in priority:
// explicit -> GITHUB_TOKEN -> GH_TOKEN.
func ResolveToken(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		return v, nil
	}
	if v := os.Getenv("GH_TOKEN"); v != "" {
		return v, nil
	}
	return "", ErrNoToken
}

