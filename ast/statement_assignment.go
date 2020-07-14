package ast

import (
	"bytes"

	"ghostlang.org/ghost/token"
)

type AssignmentStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode() {}

func (as *AssignmentStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(as.TokenLiteral() + " ")
	out.WriteString(as.Value.String())

	return out.String()
}
