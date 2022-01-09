package ast

import "ghostlang.org/x/ghost/token"

type Class struct {
	ExpressionNode
	Token token.Token
	Name  *Identifier
	Super *Identifier
	Body  *Block
}
