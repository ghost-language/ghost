// Package source keeps the text of everything Ghost has scanned, so that a
// failure can quote the line it happened on.
//
// A token knows its file, line, and column, but by the time an error is built
// the source itself is long gone — the evaluator holds a tree, not text. Rather
// than thread the source through every node, the scanner files a copy here as
// it starts, and a report asks for the line it needs when it needs it. Files
// are few and small, and this is only read on the failing path.
package source

import (
	"strings"
	"sync"
)

var (
	mutex sync.RWMutex
	files = map[string][]string{}
)

// Register records the text of a file under the name its tokens will carry.
// Registering the same name again replaces what was there, which is what the
// REPL needs: each entry is a new "repl.gs".
func Register(file string, text string) {
	if file == "" {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	files[file] = split(text)
}

// Line returns the text of a one-based line, without its terminator. It reports
// false when the file was never registered or the line is out of range, and a
// report that cannot quote a line simply does not quote one.
func Line(file string, line int) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	lines, ok := files[file]

	if !ok || line < 1 || line > len(lines) {
		return "", false
	}

	return lines[line-1], true
}

// Forget drops a file. Nothing in Ghost needs this today, but an embedder that
// runs many scripts in one process should not accumulate their sources.
func Forget(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(files, file)
}

// Reset drops every registered file.
func Reset() {
	mutex.Lock()
	defer mutex.Unlock()

	files = map[string][]string{}
}

// split breaks text into lines, tolerating both line endings. A trailing
// newline does not create a final empty line, because there is no line there to
// quote.
func split(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimSuffix(text, "\n")

	return strings.Split(text, "\n")
}
