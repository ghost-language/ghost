package token

import "fmt"

// Type is the type of the given token as a string.
type Type string

// Token contains the lexeme read by the scanner.
type Token struct {
	Type    Type        // Token type
	Lexeme  string      // String representation of literal value
	Literal interface{} // Native value in Go
	Line    int         // Line of occurance
}

func (token *Token) String() string {
	return fmt.Sprintf("%s \"%s\" %v on line %d", token.Type, token.Lexeme, token.Literal, token.Line)
}

const (
	// single-character tokens
	COLON      = ":"
	COMMA      = ","
	DOT        = "."
	LEFTBRACE  = "{"
	LEFTPAREN  = "("
	MINUS      = "-"
	PLUS       = "+"
	QUESTION   = "?"
	RIGHTBRACE = "}"
	RIGHTPAREN = ")"
	SEMICOLON  = ";"
	SLASH      = "/"
	STAR       = "*"

	// one or two character tokens
	BANG         = "!"
	BANGEQUAL    = "!="
	EQUAL        = "="
	EQUALEQUAL   = "=="
	ASSIGN       = ":="
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
	FOR      = "for"
	FUNCTION = "function"
	IF       = "if"
	NULL     = "null"
	OR       = "or"
	PRINT    = "print"
	RETURN   = "return"
	SUPER    = "super"
	THIS     = "this"
	TRUE     = "true"
	WHILE    = "while"
	EOF      = "eof"
	INVALID  = "__INVALID__"
)
