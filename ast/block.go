package ast

import "ghostlang.org/x/ghost/token"

type Block struct {
	StatementNode
	Token      token.Token
	Statements []StatementNode
}

func (node *Block) Accept(v Visitor) {
	v.visitBlock(node)
}
