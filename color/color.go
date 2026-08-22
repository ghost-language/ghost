// Package color decides whether terminal output should carry ANSI styling, and
// gives each part of a message a name instead of an escape code.
//
// Nothing outside this package should write an escape sequence. A caller asks
// for the role a piece of text plays — it is a heading, a location, a hint —
// and the profile it is writing to decides how, or whether, to dress it. That
// way a message rendered for a pipe, a log file, or a test is the same message
// as the one rendered for a terminal, minus the paint.
package color

import (
	"io"
	"os"
	"strings"
)

// Profile records whether a destination can render ANSI styling. There are only
// two of them because there is only one question worth asking: does this reader
// have a terminal, or not.
type Profile bool

const (
	// Plain writes text with no escape sequences at all.
	Plain Profile = false

	// Colored writes text wrapped in ANSI escape sequences.
	Colored Profile = true
)

// ANSI select graphic rendition sequences. These are the only escape codes in
// the codebase; every other package reaches them through a role below.
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	boldRed   = "\033[1;31m"
	boldGreen = "\033[1;32m"
	boldWhite = "\033[1;37m"
	boldCyan  = "\033[1;36m"
	boldBlue  = "\033[1;34m"
)

// Detect reports the profile appropriate for a writer.
//
// Emitting escapes into something that cannot render them is worse than
// emitting none: the reader gets `\033[1;31m` wrapped around the very message
// that was supposed to be easy to read. So the answer is only Colored once
// there is positive evidence that it will be understood, and every uncertain
// case falls back to Plain.
//
// The rules, in the order they override each other:
//
//   - NO_COLOR, set to anything at all, turns styling off. It is the reader's
//     standing instruction and nothing overrides it.
//   - A TERM of "dumb" is a terminal saying it cannot render escapes. Nothing
//     overrides that either, including the force variables below: forcing
//     colour onto a terminal that has told us it cannot show it produces
//     garbage, not colour.
//   - FORCE_COLOR and CLICOLOR_FORCE turn styling on even when the destination
//     is not a terminal at all, which is how a CI log ends up with colour.
//   - CLICOLOR=0 turns it off for a destination that would otherwise get it.
//   - Otherwise the destination has to be a terminal, and that terminal has to
//     be able to render ANSI — which on Windows means asking it to.
func Detect(writer io.Writer) Profile {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return Plain
	}

	if os.Getenv("TERM") == "dumb" {
		return Plain
	}

	if forced("FORCE_COLOR") || forced("CLICOLOR_FORCE") {
		return Colored
	}

	if os.Getenv("CLICOLOR") == "0" {
		return Plain
	}

	// Only a file can be a terminal, so anything else — a buffer, a pipe
	// wrapper, a test recorder — is answered without a syscall.
	file, ok := writer.(*os.File)

	if !ok {
		return Plain
	}

	if !isTerminal(file) {
		return Plain
	}

	if !rendersEscapes(file) {
		return Plain
	}

	return Colored
}

// forced reports whether an environment variable is set to something other than
// an explicit "0" or "false". Programs set FORCE_COLOR=1, but they also set it
// empty to mean "yes", so presence alone is the signal.
func forced(name string) bool {
	value, ok := os.LookupEnv(name)

	if !ok {
		return false
	}

	return value != "0" && !strings.EqualFold(value, "false")
}

// paint wraps text in an escape sequence, or leaves it alone when the profile
// is plain. Empty text is never painted: a zero-width run of escape codes is
// invisible clutter that still confuses anything measuring the output.
func (profile Profile) paint(code string, text string) string {
	if profile == Plain || text == "" {
		return text
	}

	return code + text + reset
}

// =============================================================================
// Roles
//
// Each of these names a part a piece of text plays in a message, not a color.
// Changing how errors look means changing the code here once.

// Error styles the heading of a failure, the part a reader's eye should land on
// first.
func (profile Profile) Error(text string) string {
	return profile.paint(boldRed, text)
}

// Warning styles the heading of something that did not stop the program.
func (profile Profile) Warning(text string) string {
	return profile.paint(yellow, text)
}

// Success styles a confirmation.
func (profile Profile) Success(text string) string {
	return profile.paint(green, text)
}

// Debug styles diagnostic output meant for whoever is working on Ghost itself.
func (profile Profile) Debug(text string) string {
	return profile.paint(blue, text)
}

// Location styles a file, line, and column reference.
func (profile Profile) Location(text string) string {
	return profile.paint(boldBlue, text)
}

// Gutter styles the line numbers and rules framing a quoted snippet. It stays
// quiet so the source itself reads normally.
func (profile Profile) Gutter(text string) string {
	return profile.paint(cyan, text)
}

// Snippet styles the quoted line of source. It is deliberately undecorated: the
// reader is meant to read it as the code they wrote.
func (profile Profile) Snippet(text string) string {
	return text
}

// Caret styles the marker pointing at the offending column.
func (profile Profile) Caret(text string) string {
	return profile.paint(boldRed, text)
}

// Help styles a suggested fix.
func (profile Profile) Help(text string) string {
	return profile.paint(boldCyan, text)
}

// Note styles a secondary remark — a call frame, a related position.
func (profile Profile) Note(text string) string {
	return profile.paint(dim, text)
}

// Emphasis styles a name or value quoted inside a sentence.
func (profile Profile) Emphasis(text string) string {
	return profile.paint(bold, text)
}
