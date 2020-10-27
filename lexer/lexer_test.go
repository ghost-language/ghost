package lexer

import (
	"testing"

	"ghostlang.org/x/ghost/token"
)

func TestNextToken(t *testing.T) {
	input := `five := 5;
ten := 10;

add := function(x, y) {
	x + y;
};

result := add(five, ten);
!*-/5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

#! \usr\bin\ghost
10 == 10;
10 != 9;

/*
This is a multiline comment
*/

"foobar"
"foo bar"
'hello world'
[1, 2]
%5
0.123
{"foo": "bar"}
while {}
true and true
true or false
1 >= 1
1 <= 1
// This is a single line comment
foo.bar
index++
index--
index += 10
index -= 10
index *= 10
index /= 10
1 .. 10
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFIER, "five"},
		{token.BIND, ":="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "ten"},
		{token.BIND, ":="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "add"},
		{token.BIND, ":="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "result"},
		{token.BIND, ":="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.ASTERISK, "*"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.GT, ">"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.NUMBER, "10"},
		{token.EQ, "=="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "10"},
		{token.NOTEQ, "!="},
		{token.NUMBER, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.STRING, "hello world"},
		{token.LBRACKET, "["},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RBRACKET, "]"},
		{token.PERCENT, "%"},
		{token.NUMBER, "5"},
		{token.NUMBER, "0.123"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.WHILE, "while"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.TRUE, "true"},
		{token.AND, "and"},
		{token.TRUE, "true"},
		{token.TRUE, "true"},
		{token.OR, "or"},
		{token.FALSE, "false"},
		{token.NUMBER, "1"},
		{token.GTE, ">="},
		{token.NUMBER, "1"},
		{token.NUMBER, "1"},
		{token.LTE, "<="},
		{token.NUMBER, "1"},
		{token.IDENTIFIER, "foo"},
		{token.DOT, "."},
		{token.IDENTIFIER, "bar"},
		{token.IDENTIFIER, "index"},
		{token.PLUSPLUS, "++"},
		{token.IDENTIFIER, "index"},
		{token.MINUSMINUS, "--"},
		{token.IDENTIFIER, "index"},
		{token.PLUSASSIGN, "+="},
		{token.NUMBER, "10"},
		{token.IDENTIFIER, "index"},
		{token.MINUSASSIGN, "-="},
		{token.NUMBER, "10"},
		{token.IDENTIFIER, "index"},
		{token.ASTERISKASSIGN, "*="},
		{token.NUMBER, "10"},
		{token.IDENTIFIER, "index"},
		{token.SLASHASSIGN, "/="},
		{token.NUMBER, "10"},
		{token.NUMBER, "1"},
		{token.RANGE, ".."},
		{token.NUMBER, "10"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		token := lexer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, token.Literal)
		}
	}
}
