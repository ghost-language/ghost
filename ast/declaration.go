package ast

import "ghostlang.org/x/ghost/token"

// Declaration ...
type Declaration struct {
	StatementNode
	Name        token.Token
	Initializer ExpressionNode
}
