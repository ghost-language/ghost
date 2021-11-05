package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) groupExpression() ast.ExpressionNode {
	parser.advance()

	group := parser.parseExpression(LOWEST)

	if !parser.checkNext(token.RIGHTPAREN) {
		return nil
	}

	return group
}
