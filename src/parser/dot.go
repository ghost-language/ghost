package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) dotExpression(left ast.ExpressionNode) ast.ExpressionNode {
	log.Debug("parsing dot expression")

	currentToken := parser.currentToken
	currentPrecedence := parser.currentTokenPrecedence()

	parser.readToken()

	if parser.nextTokenIs(token.LEFTPAREN) {
		// Method
		log.Debug("parsing method")
		expression := &ast.Method{Token: currentToken, Left: left}
		expression.Method = parser.parseExpression(currentPrecedence)

		parser.readToken()

		expression.Arguments = parser.parseExpressionList(token.RIGHTPAREN)

		return expression
	}

	// Property
	log.Debug("parsing property")
	expression := &ast.Null{}

	return expression
}
