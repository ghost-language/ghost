package ast

import (
	"bytes"
	"strings"

	"ghostlang.org/x/ghost/token"
)

type FunctionLiteral struct {
	Token      token.Token
	Name       string
	Parameters []*Identifier
	Defaults   map[string]Expression
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, p := range fl.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
