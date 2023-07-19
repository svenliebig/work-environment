package progress

import (
	"fmt"
	"math"
	"strings"
	"syscall"
	"unsafe"
)

type Writer struct {
	firstRender bool
}

func (pb *Writer) Print(p *Progress) {
	maxCols := min(int(getWidth()), 100)
	pCols := maxCols - 15

	if pb.firstRender {
		pb.firstRender = false
	} else {
		fmt.Print(strings.Repeat("\b", maxCols))
	}

	fmt.Printf("%s%% [", pad(fmt.Sprintf("%d", p.Get()), 3, " "))

	x := float64(p.Get()) / 100 * float64(pCols)
	b := float64(pCols) - x

	fmt.Print(strings.Repeat("=", int(math.Round(x))))
	if p.Get() != 100 {
		fmt.Print(">")
		b -= 1
	}
	fmt.Print(strings.Repeat(" ", int(math.Round(b))))
	fmt.Print("]")
}

func min(args ...int) int {
	var min int = args[0]

	for i := 1; i < len(args); i++ {
		if args[i] < min {
			min = args[i]
		}
	}

	return min
}
func max(args ...int) int {
	var max int = args[0]

	for i := 1; i < len(args); i++ {
		if args[i] > max {
			max = args[i]
		}
	}

	return max
}

func pad(s string, i int, padChar string) string {
	return fmt.Sprintf("%s%s", strings.Repeat(padChar, max(0, i-len(s))), s)
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() uint {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col)
}
