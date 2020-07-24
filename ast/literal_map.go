package ast

import (
	"bytes"
	"strings"

	"ghostlang.org/x/ghost/token"
)

type MapLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) expressionNode() {}

func (ml *MapLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

func (ml *MapLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for key, value := range ml.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
