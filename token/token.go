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
	Column  int         // Column of occurance on line
	File    string      // File of occurance
}

func (token *Token) String() string {
	return fmt.Sprintf("%s \"%s\" %v on line %d", token.Type, token.Lexeme, token.Literal, token.Line)
}

const (
	// single-character tokens
	COLON        = ":"
	COMMA        = ","
	LEFTBRACE    = "{"
	LEFTBRACKET  = "["
	LEFTPAREN    = "("
	MINUS        = "-"
	PLUS         = "+"
	QUESTION     = "?"
	RIGHTBRACE   = "}"
	RIGHTBRACKET = "]"
	RIGHTPAREN   = ")"
	SEMICOLON    = ";"
	SLASH        = "/"
	STAR         = "*"
	PERCENT      = "%"

	// one or two character tokens
	BANG         = "!"
	BANGEQUAL    = "!="
	DOT          = "."
	DOTDOT       = ".."
	EQUAL        = "="
	EQUALEQUAL   = "=="
	GREATER      = ">"
	GREATEREQUAL = ">="
	LESS         = "<"
	LESSEQUAL    = "<="
	PLUSEQUAL    = "+="
	MINUSEQUAL   = "-="
	STAREQUAL    = "*="
	SLASHEQUAL   = "/="

	// literals
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"
	NUMBER     = "NUMBER"

	// keywords
	AND      = "and"
	AS       = "as"
	BREAK    = "break"
	CASE     = "case"
	CLASS    = "class"
	CONTINUE = "continue"
	DEFAULT  = "default"
	ELSE     = "else"
	EXTENDS  = "extends"
	FALSE    = "false"
	FOR      = "for"
	FROM     = "from"
	FUNCTION = "function"
	IF       = "if"
	IMPORT   = "import"
	IN       = "in"
	NULL     = "null"
	OR       = "or"
	PRINT    = "print"
	RETURN   = "return"
	SUPER    = "super"
	SWITCH   = "switch"
	THIS     = "this"
	TRUE     = "true"
	WHILE    = "while"
	EOF      = "eof"
	INVALID  = "__INVALID__"
)
