// Package log writes Ghost's own messages — the ones that are not program
// errors — to the terminal.
//
// Errors a program causes go through the fault package, which knows how to
// quote source. What is left is everything Ghost says on its own behalf: what
// version is running, how long a script took, that a server could not shut down
// cleanly. Those are one line each, and they come through here.
package log

import (
	"fmt"
	"io"
	"os"

	"ghostlang.org/x/ghost/color"
)

var (
	// Out is where ordinary messages go, and Err is where problems go. Keeping
	// them apart is what lets a script's output be piped somewhere while its
	// diagnostics still reach the terminal.
	Out io.Writer = os.Stdout
	Err io.Writer = os.Stderr
)

// Debug reports something only of interest to whoever is working on Ghost.
func Debug(format string, arguments ...interface{}) {
	write(Err, color.Detect(Err).Debug("debug:")+" "+fmt.Sprintf(format, arguments...))
}

// Info reports something that went as intended.
func Info(format string, arguments ...interface{}) {
	write(Out, color.Detect(Out).Success(fmt.Sprintf(format, arguments...)))
}

// Warn reports something worth knowing that did not stop the program.
func Warn(format string, arguments ...interface{}) {
	write(Err, color.Detect(Err).Warning("warning:")+" "+fmt.Sprintf(format, arguments...))
}

// Error reports a failure of Ghost itself — a file it could not open, an
// argument it could not read. A failure inside a Ghost program is a fault, and
// is reported by the fault package instead.
func Error(format string, arguments ...interface{}) {
	write(Err, color.Detect(Err).Error("error:")+" "+fmt.Sprintf(format, arguments...))
}

// Print writes a line with no styling of its own.
func Print(text string) {
	write(Out, text)
}

func write(writer io.Writer, text string) {
	fmt.Fprintln(writer, text)
}
