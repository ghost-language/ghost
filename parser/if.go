package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ifExpression() ast.ExpressionNode {
	expression := &ast.If{Token: parser.currentToken}

	if !parser.expectNextTokenIs(token.LEFTPAREN) {
		return nil
	}

	parser.readToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectNextTokenIs(token.RIGHTPAREN) {
		return nil
	}

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	expression.Consequence = parser.blockStatement()

	if parser.nextTokenIs(token.ELSE) {
		parser.readToken()

		if parser.nextTokenIs(token.IF) {
			parser.readToken()

			expression.Alternative = &ast.Block{
				Statements: []ast.StatementNode{
					&ast.Expression{
						Expression: parser.ifExpression(),
					},
				},
			}

			return expression
		}

		if !parser.expectNextTokenIs(token.LEFTBRACE) {
			return nil
		}

		expression.Alternative = parser.blockStatement()
	}

	return expression
}
