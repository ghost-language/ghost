package ast

import "ghostlang.org/x/ghost/token"

type Block struct {
	StatementNode
	Token      token.Token
	Statements []StatementNode
}
