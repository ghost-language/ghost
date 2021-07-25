package ast

import "ghostlang.org/x/ghost/token"

type Function struct {
	Token      token.Token
	Name       token.Token
	Parameters []token.Token
	Defaults   map[string]ExpressionNode
	Body       []StatementNode
}
