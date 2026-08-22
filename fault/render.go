package fault

import (
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/color"
	"ghostlang.org/x/ghost/source"
)

const (
	// tabWidth is how wide a tab is assumed to be when a snippet is quoted.
	// The caret has to line up under the right character, and the only way to
	// be sure of that is to render tabs as a known number of spaces.
	tabWidth = 4

	// windowWidth is the widest snippet worth printing. A longer line is shown
	// as a window around the fault rather than wrapped across the terminal,
	// because a caret that has scrolled off the screen points at nothing.
	windowWidth = 96

	// ellipsis marks a snippet that has been trimmed at one end or the other.
	ellipsis = "..."
)

// Render lays the fault out in full: a heading, where it happened, the line it
// happened on with the offending part underlined, the calls it came through,
// and what to do about it.
//
//	type error: cannot use `+` between number and boolean
//	  --> example.ghost:3:12
//	   |
//	 3 | total = count + true
//	   |               ^
//	   |
//	   = help: both operands of `+` have to be the same type
//
// Sections that have nothing to say are left out, so a fault with no known
// position renders as a single line and a fault with no source on hand renders
// as a heading and a location.
func (fault *Fault) Render(profile color.Profile) string {
	if fault == nil {
		return ""
	}

	var report strings.Builder

	report.WriteString(profile.Error(fault.Kind.String() + ":"))
	report.WriteString(" ")
	report.WriteString(fault.Message)

	if !fault.Position.Known() {
		fault.renderNotes(&report, profile, 1)

		return report.String()
	}

	width := len(fmt.Sprintf("%d", fault.Position.Line))
	margin := strings.Repeat(" ", width)

	report.WriteString("\n")
	report.WriteString(margin)
	report.WriteString(profile.Gutter("--> "))
	report.WriteString(profile.Location(fault.Position.String()))

	if snippet, caret, ok := fault.snippet(); ok {
		rule := profile.Gutter(margin + " |")

		report.WriteString("\n" + rule)
		report.WriteString("\n" + profile.Gutter(fmt.Sprintf("%*d |", width, fault.Position.Line)) + " " + profile.Snippet(snippet))
		// Only the carets are painted: the run of spaces that positions them is
		// not something to color, and wrapping it would confuse anything that
		// measures the line.
		carets := strings.TrimLeft(caret, " ")
		padding := caret[:len(caret)-len(carets)]

		report.WriteString("\n" + rule + " " + padding + profile.Caret(carets))
	}

	fault.renderNotes(&report, profile, width)

	return report.String()
}

// renderNotes appends the call frames and the help line, each on its own row
// beneath the snippet.
func (fault *Fault) renderNotes(report *strings.Builder, profile color.Profile, width int) {
	if len(fault.Trace) == 0 && fault.Help == "" {
		return
	}

	margin := strings.Repeat(" ", width)
	marker := profile.Gutter(margin + " =")

	if fault.Position.Known() {
		report.WriteString("\n" + profile.Gutter(margin+" |"))
	}

	for _, frame := range fault.Trace {
		report.WriteString("\n" + marker + " " + profile.Note(frame.describe()))
	}

	if fault.Hidden > 0 {
		report.WriteString("\n" + marker + " " + profile.Note(fmt.Sprintf("and %d more call%s", fault.Hidden, suffix(fault.Hidden))))
	}

	if fault.Help != "" {
		report.WriteString("\n" + marker + " " + profile.Help("help:") + " " + fault.Help)
	}
}

// describe renders a call frame as the sentence it reads as.
func (frame Frame) describe() string {
	name := frame.Name

	if name == "" {
		name = "an anonymous function"
	}

	if !frame.Position.Known() {
		return "in " + name
	}

	return fmt.Sprintf("in %s, called at %s", name, frame.Position)
}

// snippet returns the source line the fault happened on and the caret row that
// underlines the offending part of it. It reports false when the source is not
// available, which is the case for code an embedder handed to Ghost directly.
func (fault *Fault) snippet() (string, string, bool) {
	line, ok := source.Line(fault.Position.File, fault.Position.Line)

	if !ok {
		return "", "", false
	}

	runes := []rune(line)
	column := fault.Position.Column

	if column < 1 {
		column = 1
	}

	if column > len(runes)+1 {
		column = len(runes) + 1
	}

	length := fault.Position.Length

	if length < 1 {
		length = 1
	}

	if column+length > len(runes)+1 {
		length = len(runes) + 1 - column
	}

	if length < 1 {
		length = 1
	}

	// Expand tabs before measuring anything: the caret has to be placed in the
	// same coordinate space the snippet is printed in.
	before, offset := expand(runes[:column-1], 0)
	marked, markedWidth := expand(runes[column-1:min(column-1+length, len(runes))], offset)
	after, _ := expand(runes[min(column-1+length, len(runes)):], offset+markedWidth)

	if markedWidth < 1 {
		markedWidth = 1
	}

	text := before + marked + after
	caret := strings.Repeat(" ", offset) + strings.Repeat("^", markedWidth)

	text, caret = window(text, caret, offset, markedWidth)

	return strings.TrimRight(text, " "), strings.TrimRight(caret, " "), true
}

// expand renders a run of source characters with tabs turned into spaces,
// starting from a known column so that tab stops land where they should. It
// returns the text and how many columns it occupies.
func expand(runes []rune, start int) (string, int) {
	var text strings.Builder

	width := 0

	for _, character := range runes {
		if character == '\t' {
			stop := tabWidth - (start+width)%tabWidth

			text.WriteString(strings.Repeat(" ", stop))

			width += stop

			continue
		}

		text.WriteRune(character)

		width++
	}

	return text.String(), width
}

// window trims a snippet that is too wide to print, keeping the marked part in
// view and marking each trimmed end with an ellipsis.
func window(text string, caret string, offset int, marked int) (string, string) {
	runes := []rune(text)

	if len(runes) <= windowWidth {
		return text, caret
	}

	// Center the window on the marked part, then pull it back inside the line.
	start := offset + marked/2 - windowWidth/2

	if start < 0 {
		start = 0
	}

	end := start + windowWidth

	if end > len(runes) {
		end = len(runes)
		start = end - windowWidth
	}

	trimmed := string(runes[start:end])
	shift := start

	if start > 0 {
		trimmed = ellipsis + trimmed
		shift -= len(ellipsis)
	}

	if end < len(runes) {
		trimmed += ellipsis
	}

	carets := []rune(caret)

	if shift >= len(carets) {
		return trimmed, ""
	}

	if shift < 0 {
		return trimmed, strings.Repeat(" ", -shift) + caret
	}

	return trimmed, string(carets[shift:])
}

// suffix pluralizes a count of calls.
func suffix(count int) string {
	if count == 1 {
		return ""
	}

	return "s"
}

func min(left int, right int) int {
	if left < right {
		return left
	}

	return right
}
