package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

var precedences = map[token.Type]int{
	token.EQUALEQUAL:   EQUALS,
	token.BANGEQUAL:    EQUALS,
	token.LESS:         LESSGREATER,
	token.LESSEQUAL:    LESSGREATER,
	token.GREATER:      LESSGREATER,
	token.GREATEREQUAL: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.STAR:         PRODUCT,
	token.SLASH:        PRODUCT,
}

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

	parser.registerPrefix(token.IDENTIFIER, parser.identifierLiteral)
	parser.registerPrefix(token.NUMBER, parser.numberLiteral)
	parser.registerPrefix(token.NULL, parser.nullLiteral)
	parser.registerPrefix(token.TRUE, parser.booleanLiteral)
	parser.registerPrefix(token.FALSE, parser.booleanLiteral)
	parser.registerPrefix(token.STRING, parser.stringLiteral)
	parser.registerPrefix(token.BANG, parser.prefixExpression)
	parser.registerPrefix(token.MINUS, parser.prefixExpression)

	parser.registerInfix(token.PLUS, parser.infixExpression)
	parser.registerInfix(token.MINUS, parser.infixExpression)
	parser.registerInfix(token.SLASH, parser.infixExpression)
	parser.registerInfix(token.STAR, parser.infixExpression)
	parser.registerInfix(token.EQUALEQUAL, parser.infixExpression)
	parser.registerInfix(token.BANGEQUAL, parser.infixExpression)
	parser.registerInfix(token.GREATER, parser.infixExpression)
	parser.registerInfix(token.GREATEREQUAL, parser.infixExpression)
	parser.registerInfix(token.LESS, parser.infixExpression)
	parser.registerInfix(token.LESSEQUAL, parser.infixExpression)

	return parser
}

func (parser *Parser) registerPrefix(tokenType token.Type, fn prefixParserFn) {
	parser.prefixParserFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.Type, fn infixParserFn) {
	parser.infixParserFns[tokenType] = fn
}

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

func (parser *Parser) currentPrecedence() int {
	if precedence, ok := precedences[parser.current().Type]; ok {
		// log.Debug("found precedence: %s (%d)", parser.current().Type, precedence)
		return precedence
	}

	return LOWEST
}

func (parser *Parser) nextPrecedence() int {
	if parser, ok := precedences[parser.next().Type]; ok {
		return parser
	}

	return LOWEST
}
