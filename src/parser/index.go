package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) indexExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.Index{Token: parser.currentToken, Left: left}

	parser.readToken()

	expression.Index = parser.parseExpression(LOWEST)

	if !parser.expectNextType(token.RIGHTBRACKET) {
		return nil
	}

	return expression
}
