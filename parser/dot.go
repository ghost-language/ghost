package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) dotExpression(left ast.ExpressionNode) ast.ExpressionNode {
	currentToken := parser.currentToken
	currentPrecedence := parser.currentTokenPrecedence()

	parser.readToken()

	// `new` is a keyword, so the old `Foo.new()` constructor call would
	// otherwise unravel into a cascade of unrelated syntax errors.
	if parser.currentTokenIs(token.NEW) {
		name := "Class"

		if identifier, ok := left.(*ast.Identifier); ok {
			name = identifier.Value
		}

		parser.report(parser.currentToken, "`%s` is not a method", token.NEW).
			WithHelp("construct instances with `new %s()`", name)

		return nil
	}

	if parser.nextTokenIs(token.LEFTPAREN) {
		// Method
		expression := &ast.Method{Token: currentToken, Left: left}
		expression.Method = parser.parseExpression(currentPrecedence)

		parser.readToken()

		expression.Arguments = parser.parseExpressionList(token.RIGHTPAREN)

		return expression
	}

	// Property
	expression := &ast.Property{Token: currentToken, Left: left}
	expression.Property = parser.parseExpression(currentPrecedence)

	parser.previousProperty = expression
	parser.previousIndex = nil

	return expression
}
