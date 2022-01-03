package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"github.com/shopspring/decimal"
)

func (parser *Parser) numberLiteral() ast.ExpressionNode {
	number := &ast.Number{Token: parser.currentToken}

	value, err := decimal.NewFromString(parser.currentToken.Lexeme)

	if err != nil {
		log.Error("%d:__: syntax error: could not parse %q as number", parser.currentToken.Line, parser.currentToken.Lexeme)
		return nil
	}

	number.Value = value

	return number
}
