package ast

import "ghostlang.org/x/ghost/token"

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []StatementNode
}
