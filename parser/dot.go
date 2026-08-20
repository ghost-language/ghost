package parser

import (
	"fmt"

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

		parser.errors = append(parser.errors, fmt.Sprintf(
			"%d:%d: syntax error: `%s` is not a method; construct instances with `new %s()`", parser.currentToken.Line, parser.currentToken.Column, token.NEW, name,
		))

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
