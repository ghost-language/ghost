package color

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDetectHonoursTheEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected Profile
	}{
		{"a buffer is not a terminal", nil, Plain},
		{"NO_COLOR silences styling", map[string]string{"NO_COLOR": "1"}, Plain},
		{"NO_COLOR wins over FORCE_COLOR", map[string]string{"NO_COLOR": "", "FORCE_COLOR": "1"}, Plain},
		{"FORCE_COLOR styles a pipe", map[string]string{"FORCE_COLOR": "1"}, Colored},
		{"FORCE_COLOR=0 does not", map[string]string{"FORCE_COLOR": "0"}, Plain},
		{"CLICOLOR_FORCE styles a pipe", map[string]string{"CLICOLOR_FORCE": "1"}, Colored},
		{"a dumb terminal cannot be forced", map[string]string{"TERM": "dumb", "FORCE_COLOR": "1"}, Plain},
		{"CLICOLOR=0 does not stop a force", map[string]string{"CLICOLOR": "0", "FORCE_COLOR": "1"}, Colored},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnvironment(t)

			for name, value := range test.env {
				t.Setenv(name, value)
			}

			if profile := Detect(&bytes.Buffer{}); profile != test.expected {
				t.Errorf("got=%v, expected=%v", profile, test.expected)
			}
		})
	}
}

// clearEnvironment unsets every variable Detect reads, restoring them when the
// test ends, so a case says exactly what it means rather than inheriting the
// environment the suite happens to be running in.
func clearEnvironment(t *testing.T) {
	for _, name := range []string{"NO_COLOR", "FORCE_COLOR", "CLICOLOR", "CLICOLOR_FORCE", "TERM"} {
		name := name

		value, ok := os.LookupEnv(name)

		if !ok {
			continue
		}

		t.Cleanup(func() { os.Setenv(name, value) })

		os.Unsetenv(name)
	}
}

func TestPlainProfileWritesNoEscapes(t *testing.T) {
	if got := Plain.Error("boom"); got != "boom" {
		t.Errorf("got=%q, expected=%q", got, "boom")
	}
}

func TestColoredProfileWrapsAndResets(t *testing.T) {
	got := Colored.Error("boom")

	if !strings.Contains(got, "boom") {
		t.Errorf("styling lost the text: %q", got)
	}

	if !strings.HasSuffix(got, reset) {
		t.Errorf("styling did not reset: %q", got)
	}
}

// Empty text must never be painted: a run of escape codes around nothing is
// invisible to a reader and still counts against anything measuring the line.
func TestEmptyTextIsNeverPainted(t *testing.T) {
	if got := Colored.Error(""); got != "" {
		t.Errorf("got=%q, expected empty", got)
	}
}
