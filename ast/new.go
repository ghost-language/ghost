package ast

import "ghostlang.org/x/ghost/token"

// New represents a class instantiation, e.g. `new Person("kai")`. Class holds
// the expression naming the class, so `new module.Person()` works as well.
type New struct {
	ExpressionNode
	Token     token.Token
	Class     ExpressionNode
	Arguments []ExpressionNode
}
