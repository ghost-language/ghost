package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Switch struct {
	ExpressionNode
	Token token.Token    // The "switch" token
	Value ExpressionNode // The value that will be used to determine the case
	Cases []*Case        // The cases this switch statement will handle
}
