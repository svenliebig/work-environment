package cli

import "fmt"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// colorizes cli text.
//
// usage:
//
//	cli.Colorize(cli.Red, "this is red")
func Colorize(c string, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}
