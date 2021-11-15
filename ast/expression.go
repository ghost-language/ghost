package ast

type Expression struct {
	StatementNode
	Expression ExpressionNode
}

func (structure *Expression) Accept(v Visitor) {
	v.visitExpression(structure)
}
