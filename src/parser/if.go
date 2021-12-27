package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ifExpression() ast.ExpressionNode {
	expression := &ast.If{Token: parser.currentToken}

	if !parser.expectNextType(token.LEFTPAREN) {
		return nil
	}

	parser.readToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectNextType(token.RIGHTPAREN) {
		return nil
	}

	if !parser.expectNextType(token.LEFTBRACE) {
		return nil
	}

	expression.Consequence = parser.blockStatement()

	if parser.nextTokenTypeIs(token.ELSE) {
		parser.readToken()

		if !parser.expectNextType(token.LEFTBRACE) {
			return nil
		}

		expression.Alternative = parser.blockStatement()
	}

	return expression
}
