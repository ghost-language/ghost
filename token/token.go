// Package token defines constants representing the lexical tokens of the Ghost
// programming language and basic operations on tokens.
package token

// TokenType represents the type of token.
type TokenType string

// Token represents the scanned token from source.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

// The list of tokens.
const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Basic type literals
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"

	// Operators and delimiters
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"
	DOT      = "."
	HASH     = "#"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL = "NULL"
	FOR      = "FOR"
	IN       = "IN"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	AND      = "AND"
	OR       = "OR"
	EXPORT   = "EXPORT"
	IMPORT   = "IMPORT"

	EQ    = "=="
	NOTEQ = "!="

	ASSIGN         = ":="
	PLUSASSIGN     = "+="
	MINUSASSIGN    = "-="
	ASTERISKASSIGN = "*="
	SLASHASSIGN    = "/="

	PLUSPLUS   = "++"
	MINUSMINUS = "--"

	RANGE = ".."
)

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"true":     TRUE,
	"false":    FALSE,
	"null": NULL,
	"for":      FOR,
	"in":       IN,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"while":    WHILE,
	"and":      AND,
	"or":       OR,
	"export":   EXPORT,
	"import":   IMPORT,
}

// LookupIdentifier checks the `keywords` table to see whether
// the given identifier is in fact a keyword.
func LookupIdentifier(identifier string) TokenType {
	if token, ok := keywords[identifier]; ok {
		return token
	}

	return IDENTIFIER
}
