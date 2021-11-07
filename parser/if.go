package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ifExpression() ast.ExpressionNode {
	expression := &ast.If{Token: parser.peek()}

	parser.advance() // if

	if !parser.match(token.LEFTPAREN) { // (
		log.Debug("no left paren")
		return nil
	}

	expression.Condition = parser.parseExpression(LOWEST)

	parser.advance()

	if !parser.match(token.RIGHTPAREN) {
		log.Debug("no right paren")
		return nil
	}

	if !parser.match(token.LEFTBRACE) {
		log.Debug("peek: %s", parser.peek().Type)
		log.Debug("no right brace")
		return nil
	}

	expression.Consequence = parser.blockStatement()

	if parser.checkNext(token.ELSE) {
		parser.advance()
		parser.advance()

		if !parser.match(token.LEFTBRACE) {
			log.Debug("peek: %s", parser.peek().Type)
			log.Debug("no left brace (else)")
			return nil
		}

		expression.Alternative = parser.blockStatement()
	}

	return expression
}
