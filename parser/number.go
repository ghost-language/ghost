package parser

import (
	"strconv"
	"strings"

	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) numberLiteral() ast.ExpressionNode {
	number := &ast.Number{Token: parser.currentToken}
	lexeme := parser.currentToken.Lexeme

	// If the lexeme contains a decimal point or scientific notation, parse as float.
	if strings.ContainsAny(lexeme, ".e") {
		value, err := strconv.ParseFloat(lexeme, 64)

		if err != nil {
			parser.numberError(lexeme, err)

			return nil
		}

		number.FloatValue = value
		number.IsFloat = true

		return number
	}

	value, err := strconv.ParseInt(lexeme, 10, 64)

	if err != nil {
		parser.numberError(lexeme, err)

		return nil
	}

	number.IntValue = value

	return number
}

// numberError reports a literal the scanner accepted but that is not a number
// Ghost can hold. The two ways that happens want different advice, so the
// message distinguishes a value that is too large from one that is malformed.
func (parser *Parser) numberError(lexeme string, err error) {
	raised := parser.report(parser.currentToken, "`%s` is not a valid number", lexeme)

	if numeric, ok := err.(*strconv.NumError); ok && numeric.Err == strconv.ErrRange {
		raised.WithHelp("this is outside the range Ghost can hold; the largest whole number is %d", int64(^uint64(0)>>1))
	}
}
