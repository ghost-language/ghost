package ast

import "ghostlang.org/x/ghost/token"

// Logical ...
type Logical struct {
	ExpressionNode
	Left     ExpressionNode
	Operator token.Token
	Right    ExpressionNode
}
