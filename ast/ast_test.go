package ast

import (
	"testing"

	"ghostlang.org/x/ghost/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&AssignStatement{
				Token: token.Token{Type: token.ASSIGN, Literal: ":="},
				Name: &IdentifierLiteral{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &IdentifierLiteral{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "myVar := anotherVar" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
