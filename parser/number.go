package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"github.com/shopspring/decimal"
)

func (parser *Parser) numberLiteral() ast.ExpressionNode {
	number := &ast.Number{Token: parser.peek()}

	value, err := decimal.NewFromString(parser.peek().Lexeme)

	if err != nil {
		err := error.Error{
			Reason:  error.Syntax,
			Message: fmt.Sprintf("could not parse %q as number on line %d", parser.peek().Lexeme, parser.peek().Line),
		}

		log.Error(err.Reason, err.Message)
		return nil
	}

	number.Value = value

	return number
}
