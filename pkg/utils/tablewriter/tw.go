package tablewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

var (
	_ io.Writer = &TableWriter{}
)

type TableWriter struct {
	lines           []string
	maxColumnsWidth []int
}

var CliCodes = []string{
	cli.Reset,
	cli.Red,
	cli.Green,
	cli.Yellow,
	cli.Blue,
	cli.Purple,
	cli.Cyan,
	cli.Gray,
	cli.White,
}

func (w *TableWriter) Write(p []byte) (n int, err error) {
	s := string(p)
	splits := strings.Split(s, "\t")
	w.lines = append(w.lines, s)

	for i, split := range splits {
		split = cleanse(split)

		if i > len(w.maxColumnsWidth)-1 {
			w.maxColumnsWidth = append(w.maxColumnsWidth, len(split))
		} else {
			if w.maxColumnsWidth[i] < len(split) {
				w.maxColumnsWidth[i] = len(split)
			}
		}
	}

	return 0, nil
}

func cleanse(s string) string {
	for _, c := range CliCodes {
		s = strings.ReplaceAll(s, c, "")
	}
	return s
}

func (w *TableWriter) Print() {
	for _, line := range w.lines {
		for i, column := range strings.Split(line, "\t") {
			w := w.maxColumnsWidth[i] - len(cleanse(column))
			fmt.Printf("%s%s", column, strings.Repeat(" ", w))
		}
		fmt.Printf("\n")
	}
}
