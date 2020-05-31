package repl

import (
	"bufio"
	"fmt"
	"io"

	"ghostlang.org/ghost/lexer"
	"ghostlang.org/ghost/token"
)

// PROMPT designates the REPL prompt characters to accept
// user input.
const PROMPT = ">>> "

// Start will initiate a new REPL session.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)

		for currentToken := lexer.NextToken(); currentToken.Type != token.EOF; currentToken = lexer.NextToken() {
			fmt.Printf("%+v\n", currentToken)
		}
	}
}
