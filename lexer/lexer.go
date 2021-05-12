package lexer

import (
	"fmt"

	"ghostlang.org/x/ghost/builtins"
	"ghostlang.org/x/ghost/token"
)

// Lexer takes source code as input and outputs the resulting tokens.
type Lexer struct {
	input        []rune
	character    rune
	position     int
	readPosition int
	line         int
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	lexer := &Lexer{
		input: []rune(input),
		line:  1,
	}

	lexer.readCharacter()

	return lexer
}

func (lexer *Lexer) newToken(tokenType token.TokenType) token.Token {
	literal := string(lexer.character)
	line := lexer.line

	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
	}
}

func (lexer *Lexer) newTokenWithLiteral(tokenType token.TokenType, literal string) token.Token {
	line := lexer.line

	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
	}
}

func (lexer *Lexer) newTwoCharacterToken(tokenType token.TokenType) token.Token {
	character := lexer.character
	lexer.readCharacter()
	literal := string(character) + string(lexer.character)
	line := lexer.line

	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
	}
}

// NextToken looks at the current character, and returns
// a token depending on whicharacter character it is.
func (lexer *Lexer) NextToken() token.Token {
	var currentToken token.Token

	lexer.skipWhitespace()

	switch lexer.character {
	case rune('='):
		if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.EQ)
		} else {
			currentToken = lexer.newToken(token.ASSIGN)
		}
	case rune('+'):
		if lexer.peekCharacter() == rune('+') {
			currentToken = lexer.newTwoCharacterToken(token.PLUSPLUS)
		} else if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.PLUSASSIGN)
		} else {
			currentToken = lexer.newToken(token.PLUS)
		}
	case rune('-'):
		if lexer.peekCharacter() == rune('-') {
			currentToken = lexer.newTwoCharacterToken(token.MINUSMINUS)
		} else if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.MINUSASSIGN)
		} else {
			currentToken = lexer.newToken(token.MINUS)
		}
	case rune('!'):
		if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.NOTEQ)
		} else {
			currentToken = lexer.newToken(token.BANG)
		}
	case rune('#'):
		lexer.skipSingleLineComment()

		return lexer.NextToken()
	case rune('/'):
		if lexer.peekCharacter() == rune('/') {
			lexer.skipSingleLineComment()

			return lexer.NextToken()
		} else if lexer.peekCharacter() == rune('*') {
			lexer.skipMultiLineComment()

			return lexer.NextToken()
		} else if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.SLASHASSIGN)
		} else {
			currentToken = lexer.newToken(token.SLASH)
		}
	case rune('*'):
		if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.ASTERISKASSIGN)
		} else {
			currentToken = lexer.newToken(token.ASTERISK)
		}
	case rune('%'):
		currentToken = lexer.newToken(token.PERCENT)
	case rune('<'):
		if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.LTE)
		} else {
			currentToken = lexer.newToken(token.LT)
		}
	case rune('>'):
		if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.GTE)
		} else {
			currentToken = lexer.newToken(token.GT)
		}
	case rune(';'):
		currentToken = lexer.newToken(token.SEMICOLON)
	case rune(','):
		currentToken = lexer.newToken(token.COMMA)
	case rune(':'):
		if lexer.peekCharacter() == rune('=') {
			currentToken = lexer.newTwoCharacterToken(token.ASSIGN)
		} else {
			currentToken = lexer.newToken(token.COLON)
		}
	case rune('('):
		currentToken = lexer.newToken(token.LPAREN)
	case rune(')'):
		currentToken = lexer.newToken(token.RPAREN)
	case rune('{'):
		currentToken = lexer.newToken(token.LBRACE)
	case rune('}'):
		currentToken = lexer.newToken(token.RBRACE)
	case rune('['):
		currentToken = lexer.newToken(token.LBRACKET)
	case rune(']'):
		currentToken = lexer.newToken(token.RBRACKET)
	case rune('"'):
		currentToken = lexer.newTokenWithLiteral(token.STRING, lexer.readString('"'))
	case rune('\''):
		currentToken = lexer.newTokenWithLiteral(token.STRING, lexer.readString('\''))
	case rune('.'):
		if lexer.peekCharacter() == rune('.') {
			currentToken = lexer.newTwoCharacterToken(token.RANGE)
		} else {
			currentToken = lexer.newToken(token.DOT)
		}
	case 0:
		currentToken.Type = token.EOF
		currentToken.Literal = ""
	default:
		if isDigit(lexer.character) {
			return lexer.newTokenWithLiteral(token.NUMBER, lexer.readNumber())
		}

		literal := lexer.readIdentifier()
		tokenType := token.LookupIdentifier(literal)

		return lexer.newTokenWithLiteral(tokenType, literal)
	}

	lexer.readCharacter()

	return currentToken
}

func (lexer *Lexer) skipWhitespace() {
	for isWhitespace(lexer.character) {
		if lexer.character == rune('\n') {
			lexer.line++
		}

		lexer.readCharacter()
	}
}

func (lexer *Lexer) skipSingleLineComment() {
	for lexer.character != '\n' && lexer.character != 0 {
		lexer.readCharacter()
	}

	lexer.skipWhitespace()
}

func (lexer *Lexer) skipMultiLineComment() {
	endOfComment := false

	for !endOfComment {
		if lexer.character == rune('\n') {
			lexer.line++
		}

		if lexer.character == rune(0) {
			endOfComment = true
		}

		if lexer.character == rune('*') && lexer.peekCharacter() == rune('/') {
			endOfComment = true
			lexer.readCharacter()
		}

		lexer.readCharacter()
	}

	lexer.skipWhitespace()
}

func isLetter(character rune) bool {
	return rune('a') <= character && character <= rune('z') || rune('A') <= character && character <= rune('Z') || character == rune('_')
}

func isDigit(character rune) bool {
	return rune('0') <= character && character <= rune('9')
}

func isWhitespace(character rune) bool {
	return character == rune(' ') || character == rune('\t') || character == rune('\n') || character == rune('\r')
}

func isOperator(character rune) bool {
	return character == rune('+') || character == rune('-') || character == rune('*') || character == rune('/') || character == rune('%')
}

func isComparison(character rune) bool {
	return character == rune('=') || character == rune('!') || character == rune('>') || character == rune('<')
}

func isCompound(character rune) bool {
	return character == rune('.') || character == rune(',') || character == rune('\'') || character == rune('"') || character == rune(';') || character == rune(':')
}

func isBrace(character rune) bool {
	return character == rune('{') || character == rune('}')
}

func isBracket(character rune) bool {
	return character == rune('[') || character == rune(']')
}

func isParenthesis(character rune) bool {
	return character == rune('(') || character == rune(')')
}

func isEmpty(character rune) bool {
	return character == rune(0)
}

func isIdentifier(character rune) bool {
	return !isDigit(character) && !isWhitespace(character) && !isBrace(character) && !isBracket(character) && !isParenthesis(character) && !isOperator(character) && !isCompound(character) && !isComparison(character) && !isEmpty(character)
}

func (lexer *Lexer) readString(end rune) string {
	position := lexer.position + 1

	for {
		lexer.readCharacter()

		if lexer.character == end || lexer.character == rune(0) {
			break
		}
	}

	return string(lexer.input[position:lexer.position])
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position

	for isDigit(lexer.character) || lexer.character == rune('.') {
		lexer.readCharacter()
	}

	if lexer.character == rune('e') {
		lexer.readCharacter()

		if lexer.character == rune('-') {
			lexer.readCharacter()
		}

		for isDigit(lexer.character) {
			lexer.readCharacter()
		}
	}

	return string(lexer.input[position:lexer.position])
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	readPosition := lexer.readPosition
	identifier := []rune{}
	hasDotNotation := false

	for isIdentifier(lexer.character) || lexer.character == rune('.') {
		identifier = append(identifier, lexer.character)

		if lexer.character == rune('.') {
			hasDotNotation = true
		}

		lexer.readCharacter()
	}

	if hasDotNotation {
		if _, ok := builtins.BuiltinFunctions[string(identifier)]; ok {
			return fmt.Sprintf("%s", string(identifier))
		}

		index := 0
		dotIndex := 0

		// Calculate index of dot
		for index < len(identifier) {
			if identifier[index] == rune('.') {
				dotIndex = index
			}
			index++
		}

		identifier = identifier[:dotIndex]
		lexer.position = position
		lexer.readPosition = readPosition

		for dotIndex > 0 {
			lexer.readCharacter()
			dotIndex--
		}
	}

	return string(lexer.input[position:lexer.position])
}

func (lexer *Lexer) readCharacter() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.character = rune(0)
	} else {
		lexer.character = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition

	lexer.readPosition++
}

func (lexer *Lexer) peekCharacter() rune {
	if lexer.readPosition >= len(lexer.input) {
		return rune(0)
	}

	return lexer.input[lexer.readPosition]
}
