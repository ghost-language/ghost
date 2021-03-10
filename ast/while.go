package ast

// While ...
type While struct {
	StatementNode
	Condition ExpressionNode
	Body      StatementNode
}
