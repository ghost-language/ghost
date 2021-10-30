package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) assign() ast.ExpressionNode {
	if parser.check(token.IDENTIFIER) && parser.checkNext(token.ASSIGN) {
		parser.match(token.IDENTIFIER)
		name := parser.previous()
		parser.match(token.ASSIGN)

		value := parser.expression(LOWEST)

		if value == nil {
			return value
		}

		return &ast.Assign{Token: name, Value: value}
	}

	return nil
}
