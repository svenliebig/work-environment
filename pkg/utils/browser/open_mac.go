//go:build darwin
// +build darwin

package browser

import (
	"os/exec"
)

func Open(url string) error {
	args := []string{url}
	return exec.Command("open", args...).Start()
}
