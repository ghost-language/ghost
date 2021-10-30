package ast

import "ghostlang.org/x/ghost/token"

type Boolean struct {
	Token token.Token
	Value bool
}

func (null *Boolean) expressionNode() {}
