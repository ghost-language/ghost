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
for (key, value in data) {
	key
}
こんにちは
世界
123e4
123e-4
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{token.IDENTIFIER, "five", 1},
		{token.ASSIGN, ":=", 1},
		{token.NUMBER, "5", 1},
		{token.SEMICOLON, ";", 1},
		{token.IDENTIFIER, "ten", 2},
		{token.ASSIGN, ":=", 2},
		{token.NUMBER, "10", 2},
		{token.SEMICOLON, ";", 2},
		{token.IDENTIFIER, "add", 4},
		{token.ASSIGN, ":=", 4},
		{token.FUNCTION, "function", 4},
		{token.LPAREN, "(", 4},
		{token.IDENTIFIER, "x", 4},
		{token.COMMA, ",", 4},
		{token.IDENTIFIER, "y", 4},
		{token.RPAREN, ")", 4},
		{token.LBRACE, "{", 4},
		{token.IDENTIFIER, "x", 5},
		{token.PLUS, "+", 5},
		{token.IDENTIFIER, "y", 5},
		{token.SEMICOLON, ";", 5},
		{token.RBRACE, "}", 6},
		{token.SEMICOLON, ";", 6},
		{token.IDENTIFIER, "result", 8},
		{token.ASSIGN, ":=", 8},
		{token.IDENTIFIER, "add", 8},
		{token.LPAREN, "(", 8},
		{token.IDENTIFIER, "five", 8},
		{token.COMMA, ",", 8},
		{token.IDENTIFIER, "ten", 8},
		{token.RPAREN, ")", 8},
		{token.SEMICOLON, ";", 8},
		{token.BANG, "!", 9},
		{token.ASTERISK, "*", 9},
		{token.MINUS, "-", 9},
		{token.SLASH, "/", 9},
		{token.NUMBER, "5", 9},
		{token.SEMICOLON, ";", 9},
		{token.NUMBER, "5", 10},
		{token.LT, "<", 10},
		{token.NUMBER, "10", 10},
		{token.GT, ">", 10},
		{token.NUMBER, "5", 10},
		{token.SEMICOLON, ";", 10},
		{token.IF, "if", 12},
		{token.LPAREN, "(", 12},
		{token.NUMBER, "5", 12},
		{token.LT, "<", 12},
		{token.NUMBER, "10", 12},
		{token.RPAREN, ")", 12},
		{token.LBRACE, "{", 12},
		{token.RETURN, "return", 13},
		{token.TRUE, "true", 13},
		{token.SEMICOLON, ";", 13},
		{token.RBRACE, "}", 14},
		{token.ELSE, "else", 14},
		{token.LBRACE, "{", 14},
		{token.RETURN, "return", 15},
		{token.FALSE, "false", 15},
		{token.SEMICOLON, ";", 15},
		{token.RBRACE, "}", 16},
		{token.NUMBER, "10", 19},
		{token.EQ, "==", 19},
		{token.NUMBER, "10", 19},
		{token.SEMICOLON, ";", 19},
		{token.NUMBER, "10", 20},
		{token.NOTEQ, "!=", 20},
		{token.NUMBER, "9", 20},
		{token.SEMICOLON, ";", 20},
		{token.STRING, "foobar", 26},
		{token.STRING, "foo bar", 27},
		{token.STRING, "hello world", 28},
		{token.LBRACKET, "[", 29},
		{token.NUMBER, "1", 29},
		{token.COMMA, ",", 29},
		{token.NUMBER, "2", 29},
		{token.RBRACKET, "]", 29},
		{token.PERCENT, "%", 30},
		{token.NUMBER, "5", 30},
		{token.NUMBER, "0.123", 31},
		{token.LBRACE, "{", 32},
		{token.STRING, "foo", 32},
		{token.COLON, ":", 32},
		{token.STRING, "bar", 32},
		{token.RBRACE, "}", 32},
		{token.WHILE, "while", 33},
		{token.LBRACE, "{", 33},
		{token.RBRACE, "}", 33},
		{token.TRUE, "true", 34},
		{token.AND, "and", 34},
		{token.TRUE, "true", 34},
		{token.TRUE, "true", 35},
		{token.OR, "or", 35},
		{token.FALSE, "false", 35},
		{token.NUMBER, "1", 36},
		{token.GTE, ">=", 36},
		{token.NUMBER, "1", 36},
		{token.NUMBER, "1", 37},
		{token.LTE, "<=", 37},
		{token.NUMBER, "1", 37},
		{token.IDENTIFIER, "foo", 39},
		{token.DOT, ".", 39},
		{token.IDENTIFIER, "bar", 39},
		{token.IDENTIFIER, "index", 40},
		{token.PLUSPLUS, "++", 40},
		{token.IDENTIFIER, "index", 41},
		{token.MINUSMINUS, "--", 41},
		{token.IDENTIFIER, "index", 42},
		{token.PLUSASSIGN, "+=", 42},
		{token.NUMBER, "10", 42},
		{token.IDENTIFIER, "index", 43},
		{token.MINUSASSIGN, "-=", 43},
		{token.NUMBER, "10", 43},
		{token.IDENTIFIER, "index", 44},
		{token.ASTERISKASSIGN, "*=", 44},
		{token.NUMBER, "10", 44},
		{token.IDENTIFIER, "index", 45},
		{token.SLASHASSIGN, "/=", 45},
		{token.NUMBER, "10", 45},
		{token.NUMBER, "1", 46},
		{token.RANGE, "..", 46},
		{token.NUMBER, "10", 46},
		{token.FOR, "for", 47},
		{token.LPAREN, "(", 47},
		{token.IDENTIFIER, "key", 47},
		{token.COMMA, ",", 47},
		{token.IDENTIFIER, "value", 47},
		{token.IN, "in", 47},
		{token.IDENTIFIER, "data", 47},
		{token.RPAREN, ")", 47},
		{token.LBRACE, "{", 47},
		{token.IDENTIFIER, "key", 48},
		{token.RBRACE, "}", 49},
		{token.IDENTIFIER, "こんにちは", 50},
		{token.IDENTIFIER, "世界", 51},
		{token.NUMBER, "123e4", 52},
		{token.NUMBER, "123e-4", 53},
		{token.EOF, "", 0},
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

		if token.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, token.Line)
		}
	}
}
