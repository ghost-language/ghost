package ast

// If ...
type If struct {
	StatementNode
	Condition ExpressionNode
	Then      StatementNode
	Else      StatementNode
}
