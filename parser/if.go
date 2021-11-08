package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ifExpression() ast.ExpressionNode {
	expression := &ast.If{Token: parser.peek()}

	// Consume IF token
	parser.advance()

	parser.consume(token.LEFTPAREN, "Expect '(' after 'if'. got=%s", parser.peek().Type)
	expression.Condition = parser.parseExpression(LOWEST)
	expression.Consequence = parser.blockStatement()

	parser.consume(token.RIGHTBRACE, "Expect '}' after if consequence. got=%s", parser.peek().Type)

	if parser.match(token.ELSE) {
		expression.Alternative = parser.blockStatement()
	}

	return expression
}
