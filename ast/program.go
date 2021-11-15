package ast

type Program struct {
	Statements []StatementNode
}

func (structure *Program) Accept(v Visitor) {
	v.visitProgram(structure)
}
