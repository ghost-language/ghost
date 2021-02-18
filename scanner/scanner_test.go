package scanner

import (
	"testing"

	"ghostlang.org/x/ghost/token"
)

func TestScanTokens(t *testing.T) {
	input := `(){},.-+;*"hello world"`
	tests := []struct {
		expectedType   token.Type
		expectedLexeme string
	}{
		{token.LEFTPAREN, "("},
		{token.RIGHTPAREN, ")"},
		{token.LEFTBRACE, "{"},
		{token.RIGHTBRACE, "}"},
		{token.COMMA, ","},
		{token.DOT, "."},
		{token.MINUS, "-"},
		{token.PLUS, "+"},
		{token.SEMICOLON, ";"},
		{token.STAR, "*"},
		{token.STRING, "\"hello world\""},
	}

	scanner := New(input)
	tokens := scanner.ScanTokens()

	if len(tests) != len(tokens)-1 {
		t.Fatalf("test - number of tokens is wrong. expected=%d, got=%d", len(tests), len(tokens)-1)
	}

	for i, test := range tests {
		if test.expectedType != tokens[i].Type {
			t.Fatalf("test[%d] - token type is wrong. expected=%q, got=%q", i, test.expectedType, tokens[i].Type)
		}

		if test.expectedLexeme != tokens[i].Lexeme {
			t.Fatalf("test[%d] - token literal is wrong. expected=%q, got=%q", i, test.expectedLexeme, tokens[i].Lexeme)
		}
	}

	if tokens[len(tokens)-1].Type != token.EOF {
		t.Fatalf("test[%d] - last token is not EOF", len(tokens)-1)
	}
}
