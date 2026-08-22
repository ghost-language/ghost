//go:build !windows

package color

import (
	"os"
	"path/filepath"
	"testing"
)

// A character device is not a terminal. /dev/null is the case that matters —
// it is what a redirect that discards output looks like — and the check this
// replaced treated it as a terminal and wrote escapes into it.
func TestIsTerminalRejectsOtherDevices(t *testing.T) {
	device, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)

	if err != nil {
		t.Skipf("cannot open %s: %s", os.DevNull, err)
	}

	defer device.Close()

	info, err := device.Stat()

	if err != nil {
		t.Fatalf("cannot stat %s: %s", os.DevNull, err)
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		t.Skipf("%s is not a character device here, so it does not test anything", os.DevNull)
	}

	if isTerminal(device) {
		t.Errorf("%s was taken for a terminal", os.DevNull)
	}
}

func TestIsTerminalRejectsFilesAndPipes(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "report.txt"))

	if err != nil {
		t.Fatalf("cannot create a file: %s", err)
	}

	defer file.Close()

	if isTerminal(file) {
		t.Error("a regular file was taken for a terminal")
	}

	reader, writer, err := os.Pipe()

	if err != nil {
		t.Fatalf("cannot create a pipe: %s", err)
	}

	defer reader.Close()
	defer writer.Close()

	if isTerminal(writer) {
		t.Error("a pipe was taken for a terminal")
	}
}

// The positive case has to be real too, or the detection could reject
// everything and still pass. A pty master is a terminal by any definition.
func TestIsTerminalAcceptsATerminal(t *testing.T) {
	terminal := openTerminal(t)

	defer terminal.Close()

	if !isTerminal(terminal) {
		t.Error("a pty was not recognised as a terminal")
	}
}

// Detect has to agree with isTerminal end to end: a terminal that announces
// itself through TERM gets colour, and the same terminal gets none once TERM
// says it cannot show it.
func TestDetectOnARealTerminal(t *testing.T) {
	terminal := openTerminal(t)

	defer terminal.Close()

	tests := []struct {
		name     string
		env      map[string]string
		expected Profile
	}{
		{"a terminal that announces itself", map[string]string{"TERM": "xterm-256color"}, Colored},
		{"a terminal that says it is dumb", map[string]string{"TERM": "dumb"}, Plain},
		{"a terminal that says nothing at all", nil, Plain},
		{"a reader who does not want colour", map[string]string{"TERM": "xterm", "NO_COLOR": "1"}, Plain},
		{"a reader who turned it off", map[string]string{"TERM": "xterm", "CLICOLOR": "0"}, Plain},
		{"a reader who turned it back on", map[string]string{"TERM": "xterm", "CLICOLOR": "1"}, Colored},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnvironment(t)

			for name, value := range test.env {
				t.Setenv(name, value)
			}

			if profile := Detect(terminal); profile != test.expected {
				t.Errorf("got=%v, expected=%v", profile, test.expected)
			}
		})
	}
}

// openTerminal opens a pty master, which is a terminal without needing one to
// be attached to the test run.
func openTerminal(t *testing.T) *os.File {
	t.Helper()

	terminal, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)

	if err != nil {
		t.Skipf("no pty available here: %s", err)
	}

	return terminal
}
