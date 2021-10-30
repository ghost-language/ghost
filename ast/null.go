package ast

import "ghostlang.org/x/ghost/token"

type Null struct {
	Token token.Token
}

func (null *Null) expressionNode() {}
