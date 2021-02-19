package ast

import "go/token"

// Binary structures are for binary expressions.
type Binary struct {
	Expression
	Left     Expression
	Operator token.Token
	Right    Expression
}
