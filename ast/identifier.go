package ast

import "ghostlang.org/x/ghost/token"

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode() {}
