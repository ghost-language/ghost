package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) dotExpression(left ast.ExpressionNode) ast.ExpressionNode {
	currentToken := parser.currentToken
	currentPrecedence := parser.currentTokenPrecedence()

	parser.readToken()

	if parser.nextTokenIs(token.LEFTPAREN) {
		// Method
		expression := &ast.Method{Token: currentToken, Left: left}
		expression.Method = parser.parseExpression(currentPrecedence)

		parser.readToken()

		expression.Arguments = parser.parseExpressionList(token.RIGHTPAREN)

		return expression
	}

	// Property
	log.Debug("properties TBD...")
	expression := &ast.Null{}

	return expression
}
