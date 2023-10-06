package ast

import "ghostlang.org/x/ghost/token"

type ImportFrom struct {
	ExpressionNode
	Token       token.Token
	Path        *String
	Identifiers map[string]*Identifier
	Everything  bool
}
