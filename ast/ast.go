package ast

type Node interface{}

type StatementNode interface {
	Node
	statementNode()
}

type ExpressionNode interface {
	Node
	expressionNode()
}
