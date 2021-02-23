package ast

import "ghostlang.org/x/ghost/token"

// Ternary structures are for conditional operator expressions.
type Ternary struct {
	ExpressionNode
	Condition ExpressionNode
	Question  token.Token
	Then      ExpressionNode
	Colon     token.Token
	Else      ExpressionNode
}
