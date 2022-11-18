package wepath

import "os"

func Exists(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}
