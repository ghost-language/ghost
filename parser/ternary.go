package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) ternaryExpression(left ast.ExpressionNode) ast.ExpressionNode {
	if parser.inTernaryExpression {
		return nil
	}

	parser.inTernaryExpression = true

	defer func() {
		parser.inTernaryExpression = false
	}()

	expression := &ast.Ternary{
		Token:     parser.currentToken,
		Condition: left,
	}

	// read the "?" token
	parser.readToken()

	precedence := parser.currentTokenPrecedence()
	expression.IfTrue = parser.parseExpression(precedence)

	if !parser.expectNextTokenIs(token.COLON) {
		return nil
	}

	// setup the next token
	parser.readToken()
	expression.IfFalse = parser.parseExpression(precedence)

	return expression
}
