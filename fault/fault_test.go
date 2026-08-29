package fault

import (
	"strings"
	"testing"

	"ghostlang.org/x/ghost/color"
	"ghostlang.org/x/ghost/source"
	"ghostlang.org/x/ghost/token"
)

func TestStringIsTheOneLineForm(t *testing.T) {
	raised := From(Type, Position{File: "example.gs", Line: 3, Column: 9, Length: 1}, "cannot use `+` between number and string")

	expected := "example.gs:3:9: type error: cannot use `+` between number and string"

	if raised.String() != expected {
		t.Errorf("got=%q, expected=%q", raised.String(), expected)
	}
}

func TestStringWithoutAPositionOmitsIt(t *testing.T) {
	raised := New(Name, "`greet` is not defined")

	if raised.String() != "name error: `greet` is not defined" {
		t.Errorf("got=%q", raised.String())
	}
}

func TestRenderQuotesTheOffendingLine(t *testing.T) {
	source.Reset()
	source.Register("example.gs", "x = 1\ntotal = count + label\nprint(total)\n")

	raised := From(Type, Position{File: "example.gs", Line: 2, Column: 15, Length: 1}, "cannot use `+` between number and string").
		WithHelp("both sides of `+` have to be the same type")

	expected := strings.Join([]string{
		"type error: cannot use `+` between number and string",
		" --> example.gs:2:15",
		"  |",
		"2 | total = count + label",
		"  |               ^",
		"  |",
		"  = help: both sides of `+` have to be the same type",
	}, "\n")

	if got := raised.Render(color.Plain); got != expected {
		t.Errorf("got:\n%s\n\nexpected:\n%s", got, expected)
	}
}

// The caret has to be as wide as the thing it points at, so a misspelled name
// is underlined rather than pricked at its first letter.
func TestRenderUnderlinesTheWholeLexeme(t *testing.T) {
	source.Reset()
	source.Register("example.gs", "print(nmae)\n")

	raised := From(Name, Position{File: "example.gs", Line: 1, Column: 7, Length: 4}, "`nmae` is not defined")

	if !strings.Contains(raised.Render(color.Plain), "      ^^^^") {
		t.Errorf("caret does not span the name:\n%s", raised.Render(color.Plain))
	}
}

// Tabs are printed as spaces so that the caret lands under the right character;
// if the snippet and the caret disagreed the report would point at nothing.
func TestRenderExpandsTabs(t *testing.T) {
	source.Reset()
	source.Register("example.gs", "\tx = y\n")

	raised := From(Name, Position{File: "example.gs", Line: 1, Column: 6, Length: 1}, "`y` is not defined")
	report := raised.Render(color.Plain)

	line := findLine(t, report, "1 | ")
	caret := findLine(t, report, "  | ")

	if strings.Contains(line, "\t") {
		t.Errorf("tab survived into the snippet: %q", line)
	}

	if strings.Index(line, "y") != strings.Index(caret, "^") {
		t.Errorf("caret is not under the name:\n%s", report)
	}
}

func TestRenderWindowsALongLine(t *testing.T) {
	source.Reset()
	source.Register("example.gs", strings.Repeat("a", 200)+"\n")

	raised := From(Name, Position{File: "example.gs", Line: 1, Column: 150, Length: 1}, "something")
	report := raised.Render(color.Plain)

	line := findLine(t, report, "1 | ")

	if len(line) > windowWidth+len("1 | ")+2*len(ellipsis) {
		t.Errorf("snippet was not trimmed: %d characters", len(line))
	}

	if !strings.Contains(line, ellipsis) {
		t.Errorf("trimmed snippet does not say so: %q", line)
	}

	caret := findLine(t, report, "  | ")

	if !strings.Contains(caret, "^") {
		t.Errorf("caret was trimmed away:\n%s", report)
	}
}

// Source is not always on hand — an embedder can hand Ghost a string that was
// never filed — and a report that cannot quote a line still has to render.
func TestRenderWithoutSourceStillReports(t *testing.T) {
	source.Reset()

	raised := From(Value, Position{File: "gone.gs", Line: 2, Column: 1, Length: 1}, "cannot divide by zero")

	expected := "value error: cannot divide by zero\n --> gone.gs:2:1"

	if got := raised.Render(color.Plain); got != expected {
		t.Errorf("got:\n%s\n\nexpected:\n%s", got, expected)
	}
}

func TestRenderListsCallFrames(t *testing.T) {
	source.Reset()

	raised := New(Name, "`x` is not defined").
		WithFrame("greet()", token.Token{File: "example.gs", Line: 9, Column: 1, Length: 5})

	if !strings.Contains(raised.Render(color.Plain), "in greet(), called at example.gs:9:1") {
		t.Errorf("frame missing:\n%s", raised.Render(color.Plain))
	}
}

// A runaway recursion produces thousands of identical frames. The first few say
// everything; the rest are counted.
func TestTraceIsCapped(t *testing.T) {
	raised := New(Value, "call depth exceeded")

	for index := 0; index < maxFrames+5; index++ {
		raised.WithFrame("loop()", token.Token{File: "example.gs", Line: 1, Column: 1, Length: 1})
	}

	if len(raised.Trace) != maxFrames {
		t.Errorf("kept %d frames, expected %d", len(raised.Trace), maxFrames)
	}

	if !strings.Contains(raised.Render(color.Plain), "and 5 more calls") {
		t.Errorf("hidden frames not counted:\n%s", raised.Render(color.Plain))
	}
}

func TestRenderStylesWhenAsked(t *testing.T) {
	source.Reset()
	source.Register("example.gs", "x = y\n")

	raised := From(Name, Position{File: "example.gs", Line: 1, Column: 5, Length: 1}, "`y` is not defined")

	plain := raised.Render(color.Plain)
	colored := raised.Render(color.Colored)

	if strings.Contains(plain, "\033[") {
		t.Errorf("plain report carries escapes: %q", plain)
	}

	if !strings.Contains(colored, "\033[") {
		t.Errorf("colored report carries none: %q", colored)
	}

	if stripEscapes(colored) != plain {
		t.Errorf("styling changed the text:\ngot=%q\nexpected=%q", stripEscapes(colored), plain)
	}
}

func TestPositionOfMeasuresTheLexeme(t *testing.T) {
	position := PositionOf(token.Token{File: "example.gs", Line: 1, Column: 4, Lexeme: "count"})

	if position.Length != 5 {
		t.Errorf("got=%d, expected=5", position.Length)
	}
}

func findLine(t *testing.T, report string, prefix string) string {
	t.Helper()

	for _, line := range strings.Split(report, "\n") {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}

	t.Fatalf("no line starting %q in:\n%s", prefix, report)

	return ""
}

// stripEscapes removes ANSI sequences, so a styled report can be compared
// against the plain one it is supposed to be a painted copy of.
func stripEscapes(text string) string {
	var stripped strings.Builder

	for index := 0; index < len(text); index++ {
		if text[index] != '\033' {
			stripped.WriteByte(text[index])

			continue
		}

		for index < len(text) && text[index] != 'm' {
			index++
		}
	}

	return stripped.String()
}
