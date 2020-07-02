package ast

import "bytes"

// Program is the root node. All programs consist of a slice of Statement(s).
type Program struct {
	Statements []Statement
}

// TokenLiteral prints the literal value of the token associated with this node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// String returns a stringified version of the AST for debugging purposes.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
