package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ifExpression() ast.ExpressionNode {
	expression := &ast.If{Token: parser.current()}

	if !parser.match(token.LEFTPAREN) {
		return nil
	}

	parser.advance()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.match(token.RIGHTPAREN) {
		return nil
	}

	if !parser.match(token.RIGHTBRACE) {
		return nil
	}

	expression.Consequence = parser.blockStatement()

	if parser.checkNext(token.ELSE) {
		parser.advance()

		if !parser.match(token.LEFTBRACE) {
			return nil
		}

		expression.Alternative = parser.blockStatement()
	}

	return expression
}
