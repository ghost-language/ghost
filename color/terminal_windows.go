package color

import (
	"os"
	"sync"
	"syscall"
)

// enableVirtualTerminalProcessing is the console mode flag that makes a Windows
// console interpret ANSI escapes rather than print them.
const enableVirtualTerminalProcessing = 0x0004

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode = kernel32.NewProc("SetConsoleMode")

	// Turning virtual terminal processing on is a change to the console, not a
	// question about it, so it is done once per handle however many messages
	// are printed.
	enabled      = map[syscall.Handle]bool{}
	enabledMutex sync.Mutex
)

// isTerminal reports whether a file is a console.
//
// Only a console has a console mode, so being able to read one is the answer.
// A redirect to a file or a pipe fails here, which is the same conclusion the
// ioctl reaches on Unix.
//
// The one terminal this turns away that could have shown colour is an MSYS or
// Cygwin pty — Git Bash, mintty — which Windows sees as a named pipe rather
// than a console. Recognising those means matching on the pipe's name, which
// would also match a plain redirect, and getting that wrong writes escape codes
// into somebody's file. Answering no costs colour in one terminal and FORCE_COLOR
// brings it back; answering yes too eagerly costs correctness everywhere.
func isTerminal(file *os.File) bool {
	var mode uint32

	return syscall.GetConsoleMode(syscall.Handle(file.Fd()), &mode) == nil
}

// rendersEscapes reports whether a console will interpret ANSI escapes, turning
// that on if it can.
//
// Windows consoles do not interpret escapes unless a program asks them to, and
// the ones too old to be asked print them literally instead. So this is not a
// question that can be answered by looking: the only way to know is to make the
// request and see whether it was accepted. A console that refuses is a console
// that would have shown the reader `\033[1;31m`, and it gets no colour.
func rendersEscapes(file *os.File) bool {
	handle := syscall.Handle(file.Fd())

	enabledMutex.Lock()
	defer enabledMutex.Unlock()

	if answer, known := enabled[handle]; known {
		return answer
	}

	answer := enableEscapes(handle)
	enabled[handle] = answer

	return answer
}

// enableEscapes asks a console to interpret ANSI escapes, reporting whether it
// will.
func enableEscapes(handle syscall.Handle) bool {
	var mode uint32

	if syscall.GetConsoleMode(handle, &mode) != nil {
		return false
	}

	if mode&enableVirtualTerminalProcessing != 0 {
		return true
	}

	result, _, _ := setConsoleMode.Call(uintptr(handle), uintptr(mode|enableVirtualTerminalProcessing))

	return result != 0
}
