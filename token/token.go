package token

import "fmt"

// Type is the type of the token given as a string
type Type string

// All available tokens
const (
	// single-character tokens
	LEFTPAREN  = "("
	RIGHTPAREN = ")"
	LEFTBRACE  = "{"
	RIGHTBRACE = "}"
	COMMA      = ","
	DOT        = "."
	MINUS      = "-"
	PLUS       = "+"
	SEMICOLON  = ";"
	SLASH      = "/"
	STAR       = "*"

	// one or two character tokens
	BANG         = "!"
	BANGEQUAL    = "!="
	EQUAL        = "="
	EQUALEQUAL   = "=="
	GREATER      = ">"
	GREATEREQUAL = ">="
	LESS         = "<"
	LESSEQUAL    = "<="

	// literals
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"
	NUMBER     = "NUMBER"

	// keywords
	AND      = "and"
	CLASS    = "class"
	ELSE     = "else"
	FALSE    = "false"
	FUNCTION = "function"
	FOR      = "for"
	IF       = "if"
	NULL     = "null"
	OR       = "or"
	RETURN   = "return"
	SUPER    = "super"
	THIS     = "this"
	TRUE     = "true"
	WHILE    = "while"
	EOF      = "eof"
	INVALID  = "__INVALID__"
)

// Token contains the lexeme read by the scanner and records the line.
type Token struct {
	Type    Type
	Lexeme  string
	Literal interface{}
	Line    int
}

func (token *Token) String() string {
	return fmt.Sprintf("%s \"%s\" %v on line %d", token.Type, token.Lexeme, token.Literal, token.Line)
}
