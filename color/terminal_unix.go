//go:build !windows

package color

import (
	"os"
	"syscall"
	"unsafe"
)

// isTerminal reports whether a file is a terminal.
//
// It asks the file for its terminal settings, which only a terminal has. The
// obvious shortcut — checking for a character device — is not the same
// question: /dev/null is a character device, and so is a serial port, and
// neither of them is going to show anybody a colour.
func isTerminal(file *os.File) bool {
	var termios syscall.Termios

	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		file.Fd(),
		uintptr(readTermios),
		uintptr(unsafe.Pointer(&termios)),
		0,
		0,
		0,
	)

	return errno == 0
}

// rendersEscapes reports whether the terminal on the other end understands ANSI.
//
// On a Unix terminal that comes down to TERM, which is how a terminal announces
// what it can do. An unset TERM is not an assurance of anything — it is the
// absence of one — and this is the wrong place to be optimistic, so it is read
// the same way as "dumb".
func rendersEscapes(file *os.File) bool {
	term := os.Getenv("TERM")

	return term != "" && term != "dumb"
}
