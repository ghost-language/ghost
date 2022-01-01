package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) forExpression() ast.ExpressionNode {
	log.Debug("parsing for expression")
	expression := &ast.For{Token: parser.currentToken}

	if !parser.expectNextTokenIs(token.LEFTPAREN) {
		return nil
	}

	parser.readToken()

	if !parser.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	log.Debug("for 0")

	expression.Identifier = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	expression.Initializer = parser.assign()

	log.Debug("for 1")

	if expression.Initializer == nil {
		return nil
	}

	log.Debug("for 2")

	if !parser.expectNextTokenIs(token.SEMICOLON) {
		return nil
	}

	log.Debug("for 3")

	expression.Condition = parser.parseExpression(LOWEST)

	log.Debug("for 4")

	if expression.Condition == nil {
		return nil
	}

	log.Debug("for 5")

	parser.readToken()
	parser.readToken()

	expression.Increment = parser.assign()

	if expression.Increment == nil {
		return nil
	}

	if !parser.expectNextTokenIs(token.RIGHTPAREN) {
		return nil
	}

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	expression.Block = parser.blockStatement()

	log.Debug("done parsing for expression")

	return expression
}
