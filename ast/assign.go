package ast

import "ghostlang.org/x/ghost/token"

type Assign struct {
	Token token.Token
	Value ExpressionNode
}

func (assign *Assign) statementNode() {}
