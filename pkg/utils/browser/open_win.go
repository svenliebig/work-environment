//go:build windows
// +build windows

package browser

import (
	"os/exec"
)

func Open(url string) error {
	var args []string := append([]string{"/c", "start"}, url)
	return exec.Command("cmd", args...).Start()
}
