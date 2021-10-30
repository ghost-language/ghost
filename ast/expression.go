package ast

type Expression struct {
	Expression ExpressionNode
}

func (expression *Expression) statementNode() {}
