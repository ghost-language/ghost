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

	// A dotted class name swallows its own call - `new a.Foo(1)` arrives here
	// as a single method-call node, since dotExpression (parser/dot.go) folds
	// `x.y(...)` into one node regardless of what precedence asked for it -
	// and, if anything is chained after the constructor call (`new
	// a.Foo(1).bar()`), that call is nested arbitrarily deep inside further
	// method/property nodes built on top of it. splitConstructor walks back
	// down to the first call in the chain (the constructor's) and splices it
	// out, handing back whatever was chained after it (nil if nothing was)
	// for attachConstructor to graft the finished New node into.
	if rest, path, arguments, ok := splitConstructor(class); ok {
		expression.Class = path
		expression.Arguments = arguments

		if rest == nil {
			return expression
		}

		attachConstructor(rest, expression)

		return rest
	}

	expression.Class = class

	if parser.nextTokenIs(token.LEFTPAREN) {
		parser.readToken()

		expression.Arguments = parser.parseExpressionList(token.RIGHTPAREN)
	}

	return expression
}

// splitConstructor finds the first call in a dotted class-and-call chain -
// the constructor's - and splits it into a class path and arguments. It
// only ever matches a *ast.Method or *ast.Property, so a bare (undotted)
// class name is left for newExpression's own explicit argument-list parsing,
// exactly as before this existed.
//
// rest is whatever in the chain came after the constructor call, with the
// call's own node spliced out (left as a nil Left for attachConstructor to
// fill), or nil if nothing came after it. ok is false for a plain dotted
// path with no call anywhere in it at all (`new a.b.Point`, no parens) -
// that is not this function's problem to solve; newExpression's fallback
// handles a class with no arguments the same way it always has.
func splitConstructor(expr ast.ExpressionNode) (rest ast.ExpressionNode, path ast.ExpressionNode, arguments []ast.ExpressionNode, ok bool) {
	switch node := expr.(type) {
	case *ast.Method:
		if left, innerPath, innerArguments, found := splitConstructor(node.Left); found {
			node.Left = left

			return node, innerPath, innerArguments, true
		}

		// Nothing deeper in the chain was a call, so this Method is the
		// first one - the constructor call itself.
		return nil, &ast.Property{Token: node.Token, Left: node.Left, Property: node.Method}, node.Arguments, true
	case *ast.Property:
		if left, innerPath, innerArguments, found := splitConstructor(node.Left); found {
			node.Left = left

			return node, innerPath, innerArguments, true
		}

		return nil, nil, nil, false
	default:
		return nil, nil, nil, false
	}
}

// attachConstructor fills the hole splitConstructor left (a nil Left, as
// deep as the chain goes) with the finished New expression, so whatever was
// chained after the constructor call keeps applying to it.
func attachConstructor(node ast.ExpressionNode, class *ast.New) {
	switch node := node.(type) {
	case *ast.Method:
		if node.Left == nil {
			node.Left = class

			return
		}

		attachConstructor(node.Left, class)
	case *ast.Property:
		if node.Left == nil {
			node.Left = class

			return
		}

		attachConstructor(node.Left, class)
	}
}
