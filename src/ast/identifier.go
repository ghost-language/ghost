package ast

import "ghostlang.org/x/ghost/token"

type Identifier struct {
	ExpressionNode
	AssignmentNode
	Token token.Token
	Value string
}
