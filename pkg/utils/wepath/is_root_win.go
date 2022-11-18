//go:build windows
// +build windows

package wepath

func IsRoot(p string) bool {
	return p == "c:/"
}
