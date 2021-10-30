package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"github.com/shopspring/decimal"
)

func (parser *Parser) number() ast.ExpressionNode {
	number := &ast.Number{Token: parser.current()}

	value, err := decimal.NewFromString(parser.current().Lexeme)

	if err != nil {
		err := error.Error{
			Reason:  error.Syntax,
			Message: fmt.Sprintf("could not parse %q as number on line %d", parser.current().Lexeme, parser.current().Line),
		}

		log.Error(err.Reason, err.Message)
		return nil
	}

	number.Value = value

	return number
}
