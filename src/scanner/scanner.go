package scanner

import (
	"fmt"

	"ghostlang.org/x/ghost/token"
)

// Scanner transforms our source code into tokens.
type Scanner struct {
	source       []rune // raw source to be scanned
	character    rune   // current character being scanned
	position     int    // current position in source (pointing to current character)
	readPosition int    // current reading position in source (point to next character)
	line         int    // current line being scanned
	column       int    // current column being scanned
}

// keywords contains a list of all reserved keywords.
var keywords = map[string]token.Type{
	"and":      token.AND,
	"as":       token.AS,
	"class":    token.CLASS,
	"else":     token.ELSE,
	"extends":  token.EXTENDS,
	"false":    token.FALSE,
	"for":      token.FOR,
	"from":     token.FROM,
	"function": token.FUNCTION,
	"if":       token.IF,
	"import":   token.IMPORT,
	"in":       token.IN,
	"null":     token.NULL,
	"or":       token.OR,
	"return":   token.RETURN,
	"super":    token.SUPER,
	"this":     token.THIS,
	"true":     token.TRUE,
	"while":    token.WHILE,
}

// New creates a new scanner instance.
func New(source string) *Scanner {
	scanner := Scanner{source: []rune(source), line: 1, column: 1}

	scanner.readCharacter()

	return &scanner
}

// readCharacter reads the current character and advance the readPosition.
// It also checks if we've reached the end of our source.
func (scanner *Scanner) readCharacter() {
	if scanner.readPosition >= len(scanner.source) {
		scanner.character = rune(0)
	} else {
		scanner.character = scanner.source[scanner.readPosition]
	}

	scanner.position = scanner.readPosition
	scanner.readPosition++
	scanner.column++
}

// scanToken is responsible for scanning the current character and storing the
// correct token type for it. This is the heart of our scanner.
func (scanner *Scanner) ScanToken() token.Token {
	var scannedToken token.Token

	scanner.skipWhitespace()

	switch scanner.character {
	case rune('('):
		scannedToken = scanner.newToken(token.LEFTPAREN, "(", 1)
	case rune(')'):
		scannedToken = scanner.newToken(token.RIGHTPAREN, ")", 1)
	case rune('['):
		scannedToken = scanner.newToken(token.LEFTBRACKET, "[", 1)
	case rune(']'):
		scannedToken = scanner.newToken(token.RIGHTBRACKET, "]", 1)
	case rune('{'):
		scannedToken = scanner.newToken(token.LEFTBRACE, "{", 1)
	case rune('}'):
		scannedToken = scanner.newToken(token.RIGHTBRACE, "}", 1)
	case rune(','):
		scannedToken = scanner.newToken(token.COMMA, ",", 1)
	case rune('.'):
		scannedToken = scanner.newToken(token.DOT, ".", 1)
	case rune('-'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.MINUSEQUAL, "-=", 2)
		} else {
			scannedToken = scanner.newToken(token.MINUS, "-", 1)
		}
	case rune('+'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.PLUSEQUAL, "+=", 2)
		} else {
			scannedToken = scanner.newToken(token.PLUS, "+", 1)
		}
	case rune(';'):
		scannedToken = scanner.newToken(token.SEMICOLON, ";", 1)
	case rune('*'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.STAREQUAL, "*=", 2)
		} else {
			scannedToken = scanner.newToken(token.STAR, "*", 1)
		}
	case rune('%'):
		scannedToken = scanner.newToken(token.PERCENT, "%", 1)
	case rune('?'):
		scannedToken = scanner.newToken(token.QUESTION, "?", 1)
	case rune(':'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.ASSIGN, ":=", 2)
		} else {
			scannedToken = scanner.newToken(token.COLON, ":", 1)
		}
	case rune('!'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.BANGEQUAL, "!=", 2)
		} else {
			scannedToken = scanner.newToken(token.BANG, "!", 1)
		}
	case rune('='):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.EQUALEQUAL, "==", 2)
		} else {
			scannedToken = scanner.newToken(token.EQUAL, "=", 1)
		}
	case rune('<'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.LESSEQUAL, "<=", 2)
		} else {
			scannedToken = scanner.newToken(token.LESS, "<", 1)
		}
	case rune('>'):
		if scanner.match('=') {
			scannedToken = scanner.newToken(token.GREATEREQUAL, ">=", 2)
		} else {
			scannedToken = scanner.newToken(token.GREATER, ">", 1)
		}
	case rune('#'):
		scanner.skipSingleLineComment()

		return scanner.ScanToken()
	case rune('/'):
		if scanner.match('/') {
			scanner.skipSingleLineComment()

			return scanner.ScanToken()
		} else if scanner.match('*') {
			scanner.skipMultiLineComment()

			return scanner.ScanToken()
		} else if scanner.match('=') {
			scannedToken = scanner.newToken(token.SLASHEQUAL, "/=", 2)
		} else {
			scannedToken = scanner.newToken(token.SLASH, "/", 1)
		}
	case rune('"'):
		value := scanner.scanString('"')

		scannedToken = scanner.newToken(token.STRING, value, len(value))
	case rune('\''):
		value := scanner.scanString('\'')

		scannedToken = scanner.newToken(token.STRING, value, len(value))
	case 0:
		scannedToken = scanner.newToken(token.EOF, "", 1)
	default:
		if isDigit(scanner.character) {
			number := scanner.scanNumber()

			return scanner.newToken(token.NUMBER, number, len(number))
		}

		identifier := scanner.scanIdentifier()

		return scanner.newToken(lookupIdentifier(identifier), identifier, len(identifier)+1)
	}

	scanner.readCharacter()

	return scannedToken
}

// scanString consumes characters until it hits either the closing or end of
// file. If we run to the end of the file without a closing ", we report an
// error.
func (scanner *Scanner) scanString(closing rune) string {
	position := scanner.position + 1

	for {
		scanner.readCharacter()

		if scanner.character == closing || scanner.isAtEnd() {
			break
		}
	}

	return string(scanner.source[position:scanner.position])
}

// scanNumber consumes all digits for the integer part of the literal, and then
// the fractional part if we encounter a decimal point (.) followed by at least
// one digit. If we do have a fractional part, we consume all remaining digits.
func (scanner *Scanner) scanNumber() string {
	position := scanner.position

	for isDigit(scanner.character) {
		scanner.readCharacter()
	}

	// Look for a fractional part.
	if scanner.character == rune('.') && isDigit(scanner.peekCharacter()) {
		// Consume the "."
		scanner.readCharacter()

		for isDigit(scanner.character) {
			scanner.readCharacter()
		}
	}

	// Look for a scientific notion part.
	if scanner.character == rune('e') {
		// Consume the "e"
		scanner.readCharacter()

		if scanner.character == rune('-') {
			// Consume the "-"
			scanner.readCharacter()
		}

		for isDigit(scanner.character) {
			scanner.readCharacter()
		}
	}

	return string(scanner.source[position:scanner.position])
}

// scanIdentifier consumes characters until it runs out of alphanumeric
// characters.
func (scanner *Scanner) scanIdentifier() string {
	position := scanner.position

	for isIdentifier(scanner.character) {
		scanner.readCharacter()
	}

	return string(scanner.source[position:scanner.position])
}

// =============================================================================
// Helper methods

// newToken grabs the current lexeme and creates a new token for it. In this
// case, newToken is for tokens without a literal (native Go) value.
func (scanner *Scanner) newToken(tokenType token.Type, literal interface{}, length int) token.Token {
	lexeme := fmt.Sprintf("%s", literal)
	column := scanner.column - length

	return token.Token{Type: tokenType, Lexeme: lexeme, Literal: literal, Line: scanner.line, Column: column}
}

// skipSingleLineComment consumes and reads characters until it reaches the end
// of the line. Comments are lexemes but they aren't meaningful, so we simply
// discard them here.
func (scanner *Scanner) skipSingleLineComment() {
	for scanner.character != '\n' && !scanner.isAtEnd() {
		scanner.readCharacter()
	}

	scanner.skipWhitespace()
}

// skipMultiLineComment consumes and reads characters until it reaches either
// the end of our source or the closing comment delimiter (*/). Comments are
// lexemes but they aren't meaningful, so we simply discard them here.
func (scanner *Scanner) skipMultiLineComment() {
	endOfComment := false

	for !endOfComment {
		if scanner.character == rune('\n') {
			scanner.advanceLine()
		}

		if scanner.isAtEnd() {
			endOfComment = true
		}

		if scanner.character == rune('*') && scanner.match('/') {
			endOfComment = true
			scanner.readCharacter()
		}

		scanner.readCharacter()
	}

	scanner.skipWhitespace()
}

// skipWhitespace consumes and reads whitespace characters.
func (scanner *Scanner) skipWhitespace() {
	for isWhitespace(scanner.character) {
		if scanner.character == rune('\n') {
			scanner.advanceLine()
		}

		scanner.readCharacter()
	}
}

// isAtEnd tells us if we've consumed all the characters in our raw source code.
func (scanner *Scanner) isAtEnd() bool {
	return scanner.character == 0
}

// advanceLine advances the scanner's line counter and resets the column.
func (scanner *Scanner) advanceLine() {
	scanner.line++
	scanner.column = 1
}

// isDigit tells us if the passed character is a number.
func isDigit(character rune) bool {
	return rune('0') <= character && character <= rune('9')
}

// isWhitespace tells us if the passed character is a whitespace character.
func isWhitespace(character rune) bool {
	return character == rune(' ') || character == rune('\t') || character == rune('\n') || character == rune('\r')
}

// isOperator tells us if the passed character is an operator.
func isOperator(character rune) bool {
	return character == rune('+') || character == rune('-') || character == rune('*') || character == rune('/') || character == rune('%')
}

// isComparison tells us if the passed character is a comparison.
func isComparison(character rune) bool {
	return character == rune('=') || character == rune('!') || character == rune('>') || character == rune('<')
}

// isCompound tells us if the passed character is a compound.
func isCompound(character rune) bool {
	return character == rune('.') || character == rune(',') || character == rune('\'') || character == rune('"') || character == rune(';') || character == rune(':')
}

// isBrace tells us if the passed character is a brace.
func isBrace(character rune) bool {
	return character == rune('{') || character == rune('}')
}

// isBracket tells us if the passed character is a bracket.
func isBracket(character rune) bool {
	return character == rune('[') || character == rune(']')
}

// isParenthesis tells us if the passed character is a parenthesis.
func isParenthesis(character rune) bool {
	return character == rune('(') || character == rune(')')
}

// isEmpty tells us if the passed character is empty.
func isEmpty(character rune) bool {
	return character == rune(0)
}

// isIdentifier tells us if the passed character can be used in a valid identifier.
func isIdentifier(character rune) bool {
	return !isDigit(character) && !isWhitespace(character) && !isBrace(character) && !isBracket(character) && !isParenthesis(character) && !isOperator(character) && !isCompound(character) && !isComparison(character) && !isEmpty(character)
}

// lookupIdentifier looks up the string in the list of keywords and returns its
// correct token type. If not found, then we're dealing with an identifier and
// return the identifier type.
func lookupIdentifier(identifier string) token.Type {
	if token, ok := keywords[identifier]; ok {
		return token
	}

	return token.IDENTIFIER
}

// match acts as a conditional advance, only consuming the current character if
// it's what we're looking for in "expected".
func (scanner *Scanner) match(expected rune) bool {
	if scanner.isAtEnd() {
		return false
	}

	if scanner.peekCharacter() != expected {
		return false
	}

	scanner.readCharacter()

	return true
}

// peekCharacter looks at the next upcoming character. We use this to lookahead,
// useful to check for multi-character tokens.
func (scanner *Scanner) peekCharacter() rune {
	if scanner.readPosition >= len(scanner.source) {
		return rune(0)
	}

	return scanner.source[scanner.readPosition]
}
