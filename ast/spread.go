package ast

import "ghostlang.org/x/ghost/token"

// Spread is `...expr`. It is only meaningful as an element of a call's
// argument list or a list literal - see evaluator/expressions.go's
// evaluateExpressions - where it expands the list `expr` evaluates to in
// place, rather than nesting that list as a single argument or element.
// Evaluated anywhere else, it is a syntax error rather than a value.
type Spread struct {
	ExpressionNode
	Token token.Token
	Value ExpressionNode
}
