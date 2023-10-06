package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Case struct {
	ExpressionNode
	Token   token.Token      // The "case" token
	Default bool             // Is this the default branch?
	Value   []ExpressionNode // The value of the case we'll be matching against
	Body    *Block           // The block that will be evaluated if matched
}
