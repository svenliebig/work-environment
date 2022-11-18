//go:build darwin
// +build darwin

package wepath

func IsRoot(p string) bool {
	return p == "/"
}
