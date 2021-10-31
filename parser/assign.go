package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) assign() ast.StatementNode {
	if parser.check(token.IDENTIFIER) && parser.checkNext(token.ASSIGN) {
		parser.match(token.IDENTIFIER)
		name := parser.previous()
		parser.match(token.ASSIGN)

		value := parser.parseExpression(LOWEST)

		if value == nil {
			return nil
		}

		return &ast.Assign{Token: name, Value: value}
	}

	return nil
}
