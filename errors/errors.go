package errors

import (
	"fmt"

	"ghostlang.org/x/ghost/token"
)

type ErrorBag struct {
	Message string
}

// HadRuntimeError tracks if we encountered an error during runtime.
var HadRuntimeError = false

// HadParseError tracks if we encountered an error during the parsing step.
var HadParseError = false

var ParseErrorMessage string
var RuntimeErrorMessage string

// LogError ...
func LogError(line int, message string) {
	ParseErrorMessage = fmt.Sprintf("[line %v] Error: %s\n", line, message)
	HadParseError = true
}

// RuntimeError ...
func RuntimeError(message string) {
	RuntimeErrorMessage = fmt.Sprintf("%v\n", message)
	HadRuntimeError = true
}

// ParseError ...
func ParseError(t token.Token, message string) error {
	if t.Type == token.EOF {
		message = fmt.Sprintf("[line %v] Error at end of file: %s\n", t.Line, message)
	} else {
		message = fmt.Sprintf("[line %v] Error at '%s': %s\n", t.Line, t.Lexeme, message)
	}

	HadParseError = true
	ParseErrorMessage = message
	return fmt.Errorf(message)
}

func Reset() {
	HadParseError = false
	HadRuntimeError = false
	ParseErrorMessage = ""
	RuntimeErrorMessage = ""
}