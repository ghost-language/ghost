package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) groupExpression() ast.ExpressionNode {
	// Read the opening token.LEFTPAREN ("(")
	parser.readToken()

	group := parser.parseExpression(LOWEST)

	if !parser.expectNextTokenIs(token.RIGHTPAREN) {
		return nil
	}

	return group
}
