package log

import (
	"fmt"
	"os"
	"strings"

	"ghostlang.org/x/ghost/error"
)

const (
	// ANSI terminal escape codes for color output
	AnsiReset     = "\033[0;0m"
	AnsiBlue      = "\033[34;22m"
	AnsiGreen     = "\033[32;22m"
	AnsiRed       = "\033[31;22m"
	AnsiBlueBold  = "\033[34;1m"
	AnsiGreenBold = "\033[32;1m"
	AnsiRedBold   = "\033[31;1m"
)

func LogDebug(args ...string) {
	fmt.Println(AnsiBlueBold + "debug: " + AnsiBlue + strings.Join(args, " ") + AnsiReset)
}

func LogError(reason int, args ...string) {
	var state string

	switch reason {
	case error.Syntax:
		state = "syntax error"
	case error.Runtime:
		state = "runtime error"
	case error.System:
		state = "system error"
	default:
		state = "error"
	}

	fmt.Fprintln(os.Stderr, AnsiRedBold+state+": "+AnsiRed+strings.Join(args, " ")+AnsiReset)
}
