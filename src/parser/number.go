package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"github.com/shopspring/decimal"
)

func (parser *Parser) numberLiteral() ast.ExpressionNode {
	number := &ast.Number{Token: parser.currentToken}

	value, err := decimal.NewFromString(parser.currentToken.Lexeme)

	if err != nil {
		err := error.Error{
			Reason:  error.Syntax,
			Message: fmt.Sprintf("could not parse %q as number on line %d", parser.currentToken.Lexeme, parser.currentToken.Line),
		}

		log.Error(err.Reason, err.Message)
		return nil
	}

	number.Value = value

	return number
}
