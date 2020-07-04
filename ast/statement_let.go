package ast

import (
	"bytes"

	"ghostlang.org/ghost/token"
)

// LetStatement defines the node for the `let` statement which
// binds an expression to an identifier.
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral prints the literal value of the token
// associated with the let statement node.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String returns a stringified version of the AST for
// debugging purposes.
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
