package ast

import "go/token"

// Unary structures are for unary expressions.
type Unary struct {
	Expression
	Operator token.Token
	Right    Expression
}
