package ast

// Ternary structures are for conditional operator expressions.
type Ternary struct {
	ExpressionNode
	Condition ExpressionNode
	Then      ExpressionNode
	Else      ExpressionNode
}
