package ast

import "ghostlang.org/x/ghost/token"

type ImportFrom struct {
	ExpressionNode
	Token      token.Token
	Path       *String
	Identifier *Identifier
	As         *Identifier
	Value      *Expression
}
