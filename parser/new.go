package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// newExpression parses a class instantiation, e.g. `new Person("kai")`. The
// class is parsed at CALL precedence so member access binds to it while the
// argument list does not, which keeps `new Foo().bar()` reading as
// `(new Foo()).bar()`.
func (parser *Parser) newExpression() ast.ExpressionNode {
	expression := &ast.New{Token: parser.currentToken}

	parser.readToken()

	class := parser.parseExpression(CALL)

	// A dotted class name swallows its own argument list, e.g. `new a.Foo()`
	// arrives here as a method call. Unpack it back into a class and arguments.
	if method, ok := class.(*ast.Method); ok {
		expression.Class = &ast.Property{Token: method.Token, Left: method.Left, Property: method.Method}
		expression.Arguments = method.Arguments

		return expression
	}

	expression.Class = class

	if parser.nextTokenIs(token.LEFTPAREN) {
		parser.readToken()

		expression.Arguments = parser.parseExpressionList(token.RIGHTPAREN)
	}

	return expression
}
