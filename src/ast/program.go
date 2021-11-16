package ast

type Program struct {
	Statements []StatementNode
}

func (node *Program) Accept(v Visitor) {
	v.visitProgram(node)
}
