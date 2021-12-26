package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ifExpression() ast.ExpressionNode {
	expression := &ast.If{Token: parser.peek()}

	log.Debug("peek 1: %v", parser.peek().Lexeme)

	// Consume IF token
	log.Debug("consuming...")
	parser.advance()

	log.Debug("peek 2: %v", parser.peek().Lexeme)

	parser.consume(token.LEFTPAREN, "Expect '(' after 'if'. got=%s", parser.peek().Type)

	log.Debug("peek 3: %v", parser.peek().Lexeme)

	expression.Condition = parser.parseExpression(LOWEST)

	// log.Debug("peek 4: %v", parser.peek().Lexeme)
	// parser.advance()

	log.Debug("peek 4: %v", parser.peek().Lexeme)
	parser.advance()

	expression.Consequence = parser.blockStatement()

	parser.advance()

	if parser.match(token.ELSE) {
		expression.Alternative = parser.blockStatement()
	}

	return expression
}
