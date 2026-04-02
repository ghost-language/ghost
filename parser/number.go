package parser

import (
	"strconv"
	"strings"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
)

func (parser *Parser) numberLiteral() ast.ExpressionNode {
	number := &ast.Number{Token: parser.currentToken}
	lexeme := parser.currentToken.Lexeme

	// If the lexeme contains a decimal point or scientific notation, parse as float.
	if strings.ContainsAny(lexeme, ".e") {
		value, err := strconv.ParseFloat(lexeme, 64)

		if err != nil {
			log.Error("%d:__: syntax error: could not parse %q as number", parser.currentToken.Line, lexeme)
			return nil
		}

		number.FloatValue = value
		number.IsFloat = true
	} else {
		value, err := strconv.ParseInt(lexeme, 10, 64)

		if err != nil {
			log.Error("%d:__: syntax error: could not parse %q as number", parser.currentToken.Line, lexeme)
			return nil
		}

		number.IntValue = value
	}

	return number
}
