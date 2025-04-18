package tablewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

var (
	_ io.Writer = &TableWriter{}
)

func New() *TableWriter {
	return &TableWriter{}
}

type TableWriter struct {
	lines           []string
	maxColumnsWidth []int
	lastPrintedLines int // Track how many lines were printed last time
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
		splitWidth := runewidth.StringWidth(cleanse(split))

		if i > len(w.maxColumnsWidth)-1 {
			w.maxColumnsWidth = append(w.maxColumnsWidth, splitWidth)
		} else {
			if w.maxColumnsWidth[i] < splitWidth {
				w.maxColumnsWidth[i] = splitWidth
			}
		}
	}

	return len(p), nil
}

// cleanse removes ANSI color codes from a string
func cleanse(s string) string {
	for _, c := range CliCodes {
		s = strings.ReplaceAll(s, c, "")
	}
	// Also handle any other ANSI escape sequences not captured in CliCodes
	return stripANSI(s)
}

// stripANSI removes all ANSI escape sequences from a string
func stripANSI(s string) string {
	const ansiEscapeStart = '\033'
	var result strings.Builder
	inEscSeq := false

	for _, r := range s {
		if r == ansiEscapeStart {
			inEscSeq = true
			continue
		}
		
		if inEscSeq {
			if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
				inEscSeq = false
			}
			continue
		}
		
		result.WriteRune(r)
	}
	
	return result.String()
}

// adds a line of dashes to the table writer output
func (w *TableWriter) Line() {
	w.lines = append(w.lines, "$TW_LINE")
}

func (w *TableWriter) Print() {
	// Clear previous output if this isn't the first print
	if w.lastPrintedLines > 0 {
		// Move cursor up by lastPrintedLines
		fmt.Printf("\033[%dA", w.lastPrintedLines)
		// Clear lines
		for i := 0; i < w.lastPrintedLines; i++ {
			fmt.Print("\033[2K\r") // Clear current line
			if i < w.lastPrintedLines-1 {
				fmt.Print("\033[1B") // Move down one line (except for the last one)
			}
		}
		// Move cursor back to the top of the cleared section
		if w.lastPrintedLines > 1 {
			fmt.Printf("\033[%dA", w.lastPrintedLines-1)
		}
	}

	// Keep track of how many lines we're about to print
	w.lastPrintedLines = len(w.lines)

	// Print the current state
	for _, line := range w.lines {
		if line == "$TW_LINE" {
			totalWidth := 0
			for _, width := range w.maxColumnsWidth {
				totalWidth += width
			}
			fmt.Printf("%s\n", strings.Repeat("-", totalWidth))
			continue
		}

		columns := strings.Split(line, "\t")
		for i, column := range columns {
			if i >= len(w.maxColumnsWidth) {
				fmt.Print(column)
				continue
			}
			
			cleanColumn := cleanse(column)
			padding := w.maxColumnsWidth[i] - runewidth.StringWidth(cleanColumn)
			fmt.Printf("%s%s", column, strings.Repeat(" ", padding))
		}
		fmt.Printf("\n")
	}
}

func (w *TableWriter) Flush() {
	w.lines = make([]string, 0)
	w.maxColumnsWidth = make([]int, 0)
	w.lastPrintedLines = 0
}
