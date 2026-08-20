package ast

import "ghostlang.org/x/ghost/token"

// Super represents the `super` keyword, which resolves members starting from
// the superclass of the class that declared the currently executing method.
type Super struct {
	ExpressionNode
	Token token.Token
}
