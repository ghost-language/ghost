package ast

import (
	"bytes"
	"fmt"

	"ghostlang.org/x/ghost/token"
)

type ImportExpression struct {
	Token token.Token
	Name  Expression
}

func (ie *ImportExpression) expressionNode() {}

func (ie *ImportExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *ImportExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral())
	out.WriteString("(")
	out.WriteString(fmt.Sprintf("\"%s\"", ie.Name))
	out.WriteString(")")

	return out.String()
}
