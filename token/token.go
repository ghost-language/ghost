package token

import "fmt"

// Type identifies the kind of a token. It is an integer rather than a string so
// that the parser's precedence lookups and the evaluator's operator dispatch
// compare single words instead of strings, and so operator switches compile to
// jump tables. String() recovers the source spelling for error messages.
type Type int

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
	COLON Type = iota
	COMMA
	LEFTBRACE
	LEFTBRACKET
	LEFTPAREN
	MINUS
	PLUS
	QUESTION
	RIGHTBRACE
	RIGHTBRACKET
	RIGHTPAREN
	SEMICOLON
	SLASH
	STAR
	PERCENT

	// one or two character tokens
	BANG
	BANGEQUAL
	DOT
	DOTDOT
	EQUAL
	EQUALEQUAL
	GREATER
	GREATEREQUAL
	LESS
	LESSEQUAL
	PLUSEQUAL
	PLUSPLUS
	MINUSEQUAL
	MINUSMINUS
	STAREQUAL
	SLASHEQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	AS
	BREAK
	CASE
	CLASS
	CONTINUE
	DEFAULT
	ELSE
	EXTENDS
	FALSE
	FOR
	FROM
	FUNCTION
	IF
	IMPORT
	IN
	NEW
	NULL
	OR
	PRINT
	RETURN
	SUPER
	SWITCH
	THIS
	TRAIT
	TRUE
	USE
	WHILE
	EOF
	INVALID
)

// typeNames maps each token type to its source spelling. These strings are what
// appear in parser and runtime error messages.
var typeNames = [...]string{
	COLON:        ":",
	COMMA:        ",",
	LEFTBRACE:    "{",
	LEFTBRACKET:  "[",
	LEFTPAREN:    "(",
	MINUS:        "-",
	PLUS:         "+",
	QUESTION:     "?",
	RIGHTBRACE:   "}",
	RIGHTBRACKET: "]",
	RIGHTPAREN:   ")",
	SEMICOLON:    ";",
	SLASH:        "/",
	STAR:         "*",
	PERCENT:      "%",

	BANG:         "!",
	BANGEQUAL:    "!=",
	DOT:          ".",
	DOTDOT:       "..",
	EQUAL:        "=",
	EQUALEQUAL:   "==",
	GREATER:      ">",
	GREATEREQUAL: ">=",
	LESS:         "<",
	LESSEQUAL:    "<=",
	PLUSEQUAL:    "+=",
	PLUSPLUS:     "++",
	MINUSEQUAL:   "-=",
	MINUSMINUS:   "--",
	STAREQUAL:    "*=",
	SLASHEQUAL:   "/=",

	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",

	AND:      "and",
	AS:       "as",
	BREAK:    "break",
	CASE:     "case",
	CLASS:    "class",
	CONTINUE: "continue",
	DEFAULT:  "default",
	ELSE:     "else",
	EXTENDS:  "extends",
	FALSE:    "false",
	FOR:      "for",
	FROM:     "from",
	FUNCTION: "function",
	IF:       "if",
	IMPORT:   "import",
	IN:       "in",
	NEW:      "new",
	NULL:     "null",
	OR:       "or",
	PRINT:    "print",
	RETURN:   "return",
	SUPER:    "super",
	SWITCH:   "switch",
	THIS:     "this",
	TRAIT:    "trait",
	TRUE:     "true",
	USE:      "use",
	WHILE:    "while",
	EOF:      "eof",
	INVALID:  "__INVALID__",
}

// String returns the source spelling of the token type.
func (t Type) String() string {
	if int(t) < 0 || int(t) >= len(typeNames) {
		return "__INVALID__"
	}

	return typeNames[t]
}
