package token

import "testing"

func TestTokenString(t *testing.T) {
	token := Token{Type: NUMBER, Lexeme: "3", Literal: 3, Line: 40}

	if token.String() != "NUMBER \"3\" 3 on line 40" {
		t.Fatalf("expected=NUMBER \"3\" 3 on line 40, got=%q", token.String())
	}
}
