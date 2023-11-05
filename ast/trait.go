package ast

import "ghostlang.org/x/ghost/token"

type Trait struct {
	ExpressionNode
	Token token.Token
	Name  *Identifier
	Body  *Block
}
