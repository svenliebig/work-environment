package bytes

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

type FormatOptions struct {
	// Colorize the output based on the size
	//
	// - 0-5MB: Green
    //
	// - 5-10MB: Yellow
	//
	// - 10-100MB: Orange
	//
	// - 100-1000MB: Red
	//
	// - >1000MB: Purple
	Colorize bool
}

// Format returns a human-readable string for the given number of bytes
func Format(b int64, o *FormatOptions) string {
	var result string

	if b < 1024 {
		result = fmt.Sprintf("%d B", b)
	} else if b < 1024 * 1024 {
		result = fmt.Sprintf("%d KB", b / 1024)
	} else if b < 1024 * 1024 * 1024 {
		result = fmt.Sprintf("%d MB", b / (1024 * 1024))
	} else {
		result = fmt.Sprintf("%d GB", b / (1024 * 1024 * 1024))
	}

	if o != nil && o.Colorize {
		mb := b / (1024 * 1024)
		switch {
		case mb < 5:
			return cli.Colorize(cli.Green, result)
		case mb < 10:
			return cli.Colorize(cli.Yellow, result)
		case mb < 100:
			return cli.Colorize(cli.Orange, result)
		case mb < 1000:
			return cli.Colorize(cli.Red, result)
		default:
			return cli.Colorize(cli.Purple, result)
		}
	}

	return result
}

