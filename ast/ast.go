package ast

// Node is the root interface for all expressions and statements.
type Node interface{}

// ExpressionNode is the root interface for all expressions.
type ExpressionNode interface {
	Node
}

// StatementNode is the root interface for all statements.
type StatementNode interface {
	Node
}
