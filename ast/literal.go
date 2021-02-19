package ast

// Literal structures are for literal expressions.
type Literal struct {
	Expression
	Value interface{}
}
