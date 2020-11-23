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
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	lexer := &Lexer{input: []rune(input)}

	lexer.readCharacter()

	return lexer
}

func newToken(tokenType token.TokenType, character rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(character)}
}

// NextToken looks at the current character, and returns
// a token depending on whicharacter character it is.
func (lexer *Lexer) NextToken() token.Token {
	var currentToken token.Token

	lexer.skipWhitespace()

	switch lexer.character {
	case rune('='):
		if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			literal := string(character) + string(lexer.character)
			currentToken = token.Token{Type: token.EQ, Literal: literal}
		} else {
			currentToken = newToken(token.ASSIGN, lexer.character)
		}
	case rune('+'):
		if lexer.peekCharacter() == rune('+') {
			character := lexer.character
			lexer.readCharacter()
			currentToken = token.Token{Type: token.PLUSPLUS, Literal: string(character) + string(lexer.character)}
		} else if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			currentToken = token.Token{Type: token.PLUSASSIGN, Literal: string(character) + string(lexer.character)}
		} else {
			currentToken = newToken(token.PLUS, lexer.character)
		}
	case rune('-'):
		if lexer.peekCharacter() == rune('-') {
			character := lexer.character
			lexer.readCharacter()
			currentToken = token.Token{Type: token.MINUSMINUS, Literal: string(character) + string(lexer.character)}
		} else if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			currentToken = token.Token{Type: token.MINUSASSIGN, Literal: string(character) + string(lexer.character)}
		} else {
			currentToken = newToken(token.MINUS, lexer.character)
		}
	case rune('!'):
		if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			literal := string(character) + string(lexer.character)
			currentToken = token.Token{Type: token.NOTEQ, Literal: literal}
		} else {
			currentToken = newToken(token.BANG, lexer.character)
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
			character := lexer.character
			lexer.readCharacter()
			currentToken = token.Token{Type: token.SLASHASSIGN, Literal: string(character) + string(lexer.character)}
		} else {
			currentToken = newToken(token.SLASH, lexer.character)
		}
	case rune('*'):
		if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			currentToken = token.Token{Type: token.ASTERISKASSIGN, Literal: string(character) + string(lexer.character)}
		} else {
			currentToken = newToken(token.ASTERISK, lexer.character)
		}
	case rune('%'):
		currentToken = newToken(token.PERCENT, lexer.character)
	case rune('<'):
		if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			literal := string(character) + string(lexer.character)
			currentToken = token.Token{Type: token.LTE, Literal: literal}
		} else {
			currentToken = newToken(token.LT, lexer.character)
		}
	case rune('>'):
		if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			literal := string(character) + string(lexer.character)
			currentToken = token.Token{Type: token.GTE, Literal: literal}
		} else {
			currentToken = newToken(token.GT, lexer.character)
		}
	case rune(';'):
		currentToken = newToken(token.SEMICOLON, lexer.character)
	case rune(','):
		currentToken = newToken(token.COMMA, lexer.character)
	case rune(':'):
		if lexer.peekCharacter() == rune('=') {
			character := lexer.character
			lexer.readCharacter()
			literal := string(character) + string(lexer.character)
			currentToken = token.Token{Type: token.ASSIGN, Literal: literal}
		} else {
			currentToken = newToken(token.COLON, lexer.character)
		}
	case rune('('):
		currentToken = newToken(token.LPAREN, lexer.character)
	case rune(')'):
		currentToken = newToken(token.RPAREN, lexer.character)
	case rune('{'):
		currentToken = newToken(token.LBRACE, lexer.character)
	case rune('}'):
		currentToken = newToken(token.RBRACE, lexer.character)
	case rune('['):
		currentToken = newToken(token.LBRACKET, lexer.character)
	case rune(']'):
		currentToken = newToken(token.RBRACKET, lexer.character)
	case rune('"'):
		currentToken.Type = token.STRING
		currentToken.Literal = lexer.readString('"')
	case rune('\''):
		currentToken.Type = token.STRING
		currentToken.Literal = lexer.readString('\'')
	case rune('.'):
		if lexer.peekCharacter() == rune('.') {
			character := lexer.character
			lexer.readCharacter()
			literal := string(character) + string(lexer.character)
			currentToken = token.Token{Type: token.RANGE, Literal: literal}
		} else {
			currentToken = newToken(token.DOT, lexer.character)
		}
	case 0:
		currentToken.Type = token.EOF
		currentToken.Literal = ""
	default:
		if isDigit(lexer.character) {
			currentToken.Type = token.NUMBER
			currentToken.Literal = lexer.readNumber()
			return currentToken
		}

		currentToken.Literal = lexer.readIdentifier()
		currentToken.Type = token.LookupIdentifier(currentToken.Literal)

		return currentToken
	}

	lexer.readCharacter()

	return currentToken
}

func (lexer *Lexer) skipWhitespace() {
	for isWhitespace(lexer.character) {
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
		if lexer.character == 0 {
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
