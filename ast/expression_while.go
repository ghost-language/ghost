package ast

import (
	"bytes"

	"ghostlang.org/x/ghost/token"
)

type WhileExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

func (we *WhileExpression) TokenLiteral() string {
	return we.Token.Literal
}

func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Consequence.String())

	return out.String()
}
