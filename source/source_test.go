package source

import "testing"

func TestLineReturnsRegisteredText(t *testing.T) {
	Reset()
	Register("example.ghost", "first\nsecond\nthird\n")

	tests := []struct {
		line     int
		expected string
		found    bool
	}{
		{1, "first", true},
		{3, "third", true},
		{0, "", false},
		{4, "", false},
	}

	for _, test := range tests {
		line, ok := Line("example.ghost", test.line)

		if ok != test.found || line != test.expected {
			t.Errorf("line %d: got=(%q, %v), expected=(%q, %v)", test.line, line, ok, test.expected, test.found)
		}
	}
}

func TestLineReadsWindowsEndings(t *testing.T) {
	Reset()
	Register("example.ghost", "first\r\nsecond\r\n")

	if line, _ := Line("example.ghost", 2); line != "second" {
		t.Errorf("got=%q, expected=%q", line, "second")
	}
}

// The REPL files every entry under the same name, so registering again has to
// replace what was there rather than adding to it.
func TestRegisterReplaces(t *testing.T) {
	Reset()
	Register("repl.ghost", "old")
	Register("repl.ghost", "new")

	if line, _ := Line("repl.ghost", 1); line != "new" {
		t.Errorf("got=%q, expected=%q", line, "new")
	}
}

func TestUnknownFileIsNotAnError(t *testing.T) {
	Reset()

	if _, ok := Line("never-seen.ghost", 1); ok {
		t.Error("expected an unregistered file to report no line")
	}
}

func TestForgetDropsAFile(t *testing.T) {
	Reset()
	Register("example.ghost", "text")
	Forget("example.ghost")

	if _, ok := Line("example.ghost", 1); ok {
		t.Error("expected the file to have been forgotten")
	}
}
