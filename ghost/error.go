package ghost

import (
	"fmt"
	"os"
)

// HadRuntimeError tracks if we encountered an error during runtime.
var HadRuntimeError = false

// HadParseError tracks if we encountered an error during the parsing step.
var HadParseError = false

// RuntimeError ...
func RuntimeError(message string) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	HadRuntimeError = true
}

// ParseError ...
func ParseError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
	HadParseError = true
}
