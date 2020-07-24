package ast

import (
	"testing"

	"ghostlang.org/x/ghost/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
				Expression: &BindExpression{
					Token: token.Token{Type: token.BIND, Literal: ":="},
					Left: &Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}

	if program.String() != "myVar := anotherVar" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
