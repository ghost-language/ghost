package ghost

import (
	"fmt"
	"os"

	"ghostlang.org/x/ghost/token"
)

// HadRuntimeError tracks if we encountered an error during runtime.
var HadRuntimeError = false

// HadParseError tracks if we encountered an error during the parsing step.
var HadParseError = false

// Log ...
func LogError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %v] Error: %s\n", line, message)
	HadParseError = true
}

// RuntimeError ...
func RuntimeError(message string) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	HadRuntimeError = true
}

// ParseError ...
func ParseError(t token.Token, message string) error {
	message = fmt.Sprintf("[line %v] Error at '%s': %s", t.Line, t.Lexeme, message)

	if t.Type == token.EOF {
		message = fmt.Sprintf("[line %v] Error at end of file: %s", t.Line, message)
	}

	HadParseError = true
	return fmt.Errorf(message)
}
