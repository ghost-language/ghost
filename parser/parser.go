package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) Parser {
	return Parser{tokens, 0}
}

func (parser *Parser) Parse() []ast.StatementNode {
	statements := make([]ast.StatementNode, 0)

	for !parser.isAtEnd() {
		log.LogDebug(fmt.Sprintf("current: %v", parser.current))
		statement := parser.declaration()

		if statement != nil {
			statements = append(statements, statement)
		} else {
			// parser.advance()
			err := error.Error{
				Reason:  error.Runtime,
				Message: "unknown statement",
			}

			log.LogError(err.Reason, err.Message)
		}
	}

	return statements
}

// =============================================================================
// Helper methods

func (parser *Parser) match(tt ...token.Type) bool {
	for _, t := range tt {
		if parser.check(t) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) consume(tt token.Type, message string) token.Token {
	if parser.check(tt) {
		return parser.advance()
	}

	err := error.Error{
		Reason:  error.Runtime,
		Message: message,
	}

	log.LogError(err.Reason, err.Message)

	return parser.previous()
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

func (parser *Parser) checkNext(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.next().Type == tt
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) next() token.Token {
	return parser.tokens[parser.current+1]
}

func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}
