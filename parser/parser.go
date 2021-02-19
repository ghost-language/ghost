package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// Parser stores the list of tokens and current to point to the next token
// eagerly waiting to be parsed.
type Parser struct {
	tokens  []token.Token
	current int
}

// New creates a new Parser instance.
func New(tokens []token.Token) Parser {
	return Parser{tokens, 0}
}

// Parse kicks off the parser.
func (parser *Parser) Parse() ast.Expression {
	expression := parser.expression()

	return expression
}

// expression starts the process of parsing expression grammar rules.
//
// Each method for parsing a grammar rule produces a syntax tree for that rule
// and returns it to the caller. When the body of the rule contains a non-
// terminal -- a reference to another rule -- we call that other rule's method.
func (parser *Parser) expression() ast.Expression {
	return parser.equality()
}

func (parser *Parser) equality() ast.Expression {
	expression := parser.comparison()

	for parser.match(token.BANGEQUAL, token.EQUALEQUAL) {
		operator := parser.previous()
		right := parser.comparison()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

// =============================================================================
// Helper methods

// match checks if the current token has any of the given types. If so, it
// consumes the token and returns true. Otherwise, it returns false and leaves
// the current token alone.
func (parser *Parser) match(tt ...token.Type) bool {
	for _, t := range tt {
		if parser.check(t) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.current++
	}

	return parser.previous()
}

func (parser *Parser) check(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == tt
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.current]
}

// previous returns the previous token.
func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}
