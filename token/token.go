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
	Line    int         // Line the lexeme starts on
	Column  int         // Column the lexeme starts at, counting from one
	Length  int         // Width of the lexeme in characters, so a report can underline all of it
	File    string      // File of occurance
}

func (token *Token) String() string {
	return fmt.Sprintf("%s \"%s\" %v on line %d", token.Type, token.Lexeme, token.Literal, token.Line)
}

// Describe names a token the way a sentence about it needs to read. Error
// messages quote what was actually written, except at the end of the file where
// there is nothing to quote.
func (token Token) Describe() string {
	if token.Type == EOF {
		return "the end of the file"
	}

	if token.Lexeme == "" {
		return token.Type.Describe()
	}

	return "`" + token.Lexeme + "`"
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
	TEMPLATESTRING
	TEMPLATESTRINGEND

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

	IDENTIFIER:        "IDENTIFIER",
	STRING:            "STRING",
	NUMBER:            "NUMBER",
	TEMPLATESTRING:    "TEMPLATESTRING",
	TEMPLATESTRINGEND: "TEMPLATESTRINGEND",

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

// Describe names a token type the way a sentence about it needs to read. An
// operator or a keyword is quoted as it is written, because that is exactly what
// the reader has to type; a class of token, such as a name or a number, is
// described instead, because there is no single spelling to quote.
func (t Type) Describe() string {
	switch t {
	case IDENTIFIER:
		return "a name"
	case STRING:
		return "a string"
	case NUMBER:
		return "a number"
	case TEMPLATESTRING, TEMPLATESTRINGEND:
		return "a template literal"
	case EOF:
		return "the end of the file"
	case INVALID:
		return "an unreadable token"
	}

	return "`" + t.String() + "`"
}

// String returns the source spelling of the token type.
func (t Type) String() string {
	if int(t) < 0 || int(t) >= len(typeNames) {
		return "__INVALID__"
	}

	return typeNames[t]
}
