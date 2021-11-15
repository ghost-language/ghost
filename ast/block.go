package ast

import "ghostlang.org/x/ghost/token"

type Block struct {
	StatementNode
	Token      token.Token
	Statements []StatementNode
}

func (structure *Block) Accept(v Visitor) {
	v.visitBlock(structure)
}
