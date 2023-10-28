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
		`( ) [ ] { } , . - + ; * % ? : > < >= <= ! != = == "hello world" 42 3.14 6.67428e-11 foo foobar hello1 true false class whilefoo こんにちは 世界 += -= *= /= import from as .. index++ index--`,
		[]struct {
			expectedType   token.Type
			expectedLexeme string
		}{
			{token.LEFTPAREN, "("},
			{token.RIGHTPAREN, ")"},
			{token.LEFTBRACKET, "["},
			{token.RIGHTBRACKET, "]"},
			{token.LEFTBRACE, "{"},
			{token.RIGHTBRACE, "}"},
			{token.COMMA, ","},
			{token.DOT, "."},
			{token.MINUS, "-"},
			{token.PLUS, "+"},
			{token.SEMICOLON, ";"},
			{token.STAR, "*"},
			{token.PERCENT, "%"},
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
			{token.STRING, "hello world"},
			{token.NUMBER, "42"},
			{token.NUMBER, "3.14"},
			{token.NUMBER, "6.67428e-11"},
			{token.IDENTIFIER, "foo"},
			{token.IDENTIFIER, "foobar"},
			{token.IDENTIFIER, "hello1"},
			{token.TRUE, "true"},
			{token.FALSE, "false"},
			{token.CLASS, "class"},
			{token.IDENTIFIER, "whilefoo"},
			{token.IDENTIFIER, "こんにちは"},
			{token.IDENTIFIER, "世界"},
			{token.PLUSEQUAL, "+="},
			{token.MINUSEQUAL, "-="},
			{token.STAREQUAL, "*="},
			{token.SLASHEQUAL, "/="},
			{token.IMPORT, "import"},
			{token.FROM, "from"},
			{token.AS, "as"},
			{token.DOTDOT, ".."},
			{token.IDENTIFIER, "index"},
			{token.PLUSPLUS, "++"},
			{token.IDENTIFIER, "index"},
			{token.MINUSMINUS, "--"},
			{token.EOF, ""},
		},
	}

	scanner := New(test.input, "test.ghost")

	for _, tok := range test.expected {
		token := scanner.ScanToken()

		if tok.expectedType != token.Type {
			t.Fatalf("token type is wrong. expected=%q, got=%q", tok.expectedType, token.Type)
		}

		if tok.expectedLexeme != token.Lexeme {
			t.Fatalf("token lexeme is wrong. expected=%q, got=%q", tok.expectedLexeme, token.Lexeme)
		}
	}
}
