package ast

type Expression struct {
	StatementNode
	Expression ExpressionNode
}

func (node *Expression) Accept(v Visitor) {
	v.visitExpression(node)
}
