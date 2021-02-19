package ast

// Grouping structures are for parenthesised expressions.
type Grouping struct {
	ExpressionNode
	Expression ExpressionNode
}
