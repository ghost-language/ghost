package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

var precedences = map[token.Type]int{}

const (
	_ int = iota
	LOWEST
	OR
	AND
	RANGE
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	MODULO
	PREFIX
	CALL
	INDEX
)

type (
	prefixParserFn  func() ast.ExpressionNode
	infixParserFn   func(ast.ExpressionNode) ast.ExpressionNode
	postfixParserFn func() ast.ExpressionNode
)

type Parser struct {
	tokens   []token.Token
	position int

	prefixParserFns  map[token.Type]prefixParserFn
	infixParserFns   map[token.Type]infixParserFn
	postfixParserFns map[token.Type]postfixParserFn
}

func New(tokens []token.Token) *Parser {
	parser := &Parser{tokens: tokens, position: 0}

	parser.prefixParserFns = make(map[token.Type]prefixParserFn)
	parser.infixParserFns = make(map[token.Type]infixParserFn)
	parser.postfixParserFns = make(map[token.Type]postfixParserFn)

	parser.registerPrefix(token.IDENTIFIER, parser.identifier)
	parser.registerPrefix(token.NUMBER, parser.number)
	parser.registerPrefix(token.NULL, parser.null)
	parser.registerPrefix(token.TRUE, parser.boolean)
	parser.registerPrefix(token.FALSE, parser.boolean)

	return parser
}

func (parser *Parser) registerPrefix(tokenType token.Type, fn prefixParserFn) {
	parser.prefixParserFns[tokenType] = fn
}

// func (parser *Parser) registerInfix(tokenType token.Type, fn infixParserFn) {
// 	parser.infixParserFns[tokenType] = fn
// }

// func (parser *Parser) registerPostfix(tokenType token.Type, fn postfixParserFn) {
// 	parser.postfixParserFns[tokenType] = fn
// }

func (parser *Parser) Parse() []ast.StatementNode {
	statements := make([]ast.StatementNode, 0)

	for !parser.isAtEnd() {
		statement := parser.statement()

		statements = append(statements, statement)
		parser.advance()
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

// func (parser *Parser) consume(tt token.Type, message string) token.Token {
// 	if parser.check(tt) {
// 		return parser.advance()
// 	}

// 	err := error.Error{
// 		Reason:  error.Runtime,
// 		Message: message,
// 	}

// 	log.Error(err.Reason, err.Message)

// 	return parser.previous()
// }

func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.position++
	}

	return parser.previous()
}

func (parser *Parser) check(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.current().Type == tt
}

func (parser *Parser) checkNext(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.next().Type == tt
}

func (parser *Parser) isAtEnd() bool {
	return parser.current().Type == token.EOF
}

func (parser *Parser) current() token.Token {
	return parser.tokens[parser.position]
}

func (parser *Parser) next() token.Token {
	return parser.tokens[parser.position+1]
}

func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.position-1]
}

// func (parser *Parser) currentPrecedence() int {
// 	if parser, ok := precedences[parser.current().Type]; ok {
// 		return parser
// 	}

// 	return LOWEST
// }

func (parser *Parser) nextPrecedence() int {
	if parser, ok := precedences[parser.next().Type]; ok {
		return parser
	}

	return LOWEST
}
