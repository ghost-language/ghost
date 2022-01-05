package ast

import "ghostlang.org/x/ghost/token"

type Return struct {
	StatementNode
	Token token.Token
	Value ExpressionNode
}
