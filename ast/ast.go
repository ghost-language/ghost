package ast

type Node interface{}

type StatementNode interface {
	Node
}

type ExpressionNode interface {
	Node
}
