package ast

import (
	"strings"

	"ghostlang.org/x/ghost/token"
)

type Class struct {
	StatementNode
	Name token.Token
	Methods []*Function
	EnvIndex int
}

func (c *Class) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString("class")
	out.WriteString("")
	out.WriteString(c.Name.Lexeme)
	out.WriteString("")

	return out.String()
}