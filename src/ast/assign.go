package ast

import "ghostlang.org/x/ghost/token"

type Assign struct {
	StatementNode
	Token token.Token
	Name  AssignmentNode
	Value ExpressionNode
}
