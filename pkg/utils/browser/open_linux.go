//go:build linux
// +build linux

package browser

import (
	"os/exec"
)

func Open(url string) error {
	var args []string := []string{url}
	return exec.Command("xdg-open", args...).Start()
}
