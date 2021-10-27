package scanner

import (
	"testing"

	"ghostlang.org/x/ghost/token"
)

func TestScanTokens(t *testing.T) {
	test := struct {
		input    string
		expected []struct {
			expectedType   token.Type
			expectedLexeme string
		}
	}{
		`( ) { } , . - + ; * ? : > < >= <= ! != = == := "hello world" 42 3.14 6.67428e-11 foo foobar true false class whilefoo`,
		[]struct {
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
			{token.QUESTION, "?"},
			{token.COLON, ":"},
			{token.GREATER, ">"},
			{token.LESS, "<"},
			{token.GREATEREQUAL, ">="},
			{token.LESSEQUAL, "<="},
			{token.BANG, "!"},
			{token.BANGEQUAL, "!="},
			{token.EQUAL, "="},
			{token.EQUALEQUAL, "=="},
			{token.ASSIGN, ":="},
			{token.STRING, "\"hello world\""},
			{token.NUMBER, "42"},
			{token.NUMBER, "3.14"},
			{token.NUMBER, "6.67428e-11"},
			{token.IDENTIFIER, "foo"},
			{token.IDENTIFIER, "foobar"},
			{token.TRUE, "true"},
			{token.FALSE, "false"},
			{token.CLASS, "class"},
			{token.IDENTIFIER, "whilefoo"},
		},
	}

	scanner := New(test.input)
	tokens := scanner.ScanTokens()

	if len(test.expected) != len(tokens)-1 {
		t.Fatalf("number of tokens is wrong. expected=%d, got=%d", len(test.expected), len(tokens)-1)
	}

	for i, tok := range test.expected {
		if tok.expectedType != tokens[i].Type {
			t.Fatalf("token type is wrong. expected=%q, got=%q", tok.expectedType, tokens[i].Type)
		}

		if tok.expectedLexeme != tokens[i].Lexeme {
			t.Fatalf("token lexeme is wrong. expected=%q, got=%q", tok.expectedLexeme, tokens[i].Lexeme)
		}
	}

	if tokens[len(tokens)-1].Type != token.EOF {
		t.Fatalf("last token is not EOF")
	}
}
