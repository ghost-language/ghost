package ast

import (
	"bytes"

	"ghostlang.org/x/ghost/token"
)

type BindExpression struct {
	Token token.Token
	Left  Expression
	Value Expression
}

func (be *BindExpression) expressionNode() {}

func (be *BindExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BindExpression) String() string {
	var out bytes.Buffer

	out.WriteString(be.Left.String())
	out.WriteString(" " + be.TokenLiteral() + " ")
	out.WriteString(be.Value.String())

	return out.String()
}
