package cli

var (
	bold          = "\033[1m"
	italic        = "\033[2m"
	underline     = "\033[4m"
	blink         = "\033[5m"
	inverse       = "\033[7m"
	hidden        = "\033[8m"
	strikethrough = "\033[9m"
)

func Bold(s string) string {
	return bold + s + Reset
}

func Italic(s string) string {
	return italic + s + Reset
}

func Underline(s string) string {
	return underline + s + Reset
}

func Blink(s string) string {
	return blink + s + Reset
}

func Inverse(s string) string {
	return inverse + s + Reset
}

func Hidden(s string) string {
	return hidden + s + Reset
}

func Strikethrough(s string) string {
	return strikethrough + s + Reset
}
