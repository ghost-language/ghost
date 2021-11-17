package parser

import (
	"fmt"

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
	token.PERCENT:      MODULO,
	token.LEFTPAREN:    CALL,
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
	errors   []string

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
	parser.registerPrefix(token.IF, parser.ifExpression)
	parser.registerPrefix(token.LEFTPAREN, parser.groupExpression)
	parser.registerPrefix(token.FUNCTION, parser.functionStatement)

	parser.registerInfix(token.PLUS, parser.infixExpression)
	parser.registerInfix(token.MINUS, parser.infixExpression)
	parser.registerInfix(token.SLASH, parser.infixExpression)
	parser.registerInfix(token.STAR, parser.infixExpression)
	parser.registerInfix(token.PERCENT, parser.infixExpression)
	parser.registerInfix(token.EQUALEQUAL, parser.infixExpression)
	parser.registerInfix(token.BANGEQUAL, parser.infixExpression)
	parser.registerInfix(token.GREATER, parser.infixExpression)
	parser.registerInfix(token.GREATEREQUAL, parser.infixExpression)
	parser.registerInfix(token.LESS, parser.infixExpression)
	parser.registerInfix(token.LESSEQUAL, parser.infixExpression)
	parser.registerInfix(token.LEFTPAREN, parser.callExpression)

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

func (parser *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.StatementNode{}

	for !parser.isAtEnd() {
		statement := parser.statement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.advance()
	}

	return program
}

// Errors returns the slice of errors contained within the parser instance.
func (parser *Parser) Errors() []string {
	return parser.errors
}

// =============================================================================
// Helper methods

// match checks to see if the current token has any of the given types. If so,
// it consumes the token and returns true. Otherwise, it returns false and
// leaves the current token alone.
func (parser *Parser) match(tt ...token.Type) bool {
	for _, t := range tt {
		if parser.check(t) {
			parser.advance()
			return true
		}
	}

	return false
}

// advance consumes the current token and returns it, similar to how our
// scanner's corresponding method crawls through characters.
func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.position++
	}

	return parser.previous()
}

// check returns true if the current token is of the given type. Unlike match(),
// it never consumes the token, it only looks at it.
func (parser *Parser) check(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == tt
}

// consume checks to see if the next token is of the expected typ. If so, it
// consumes the token and carries on. If some other token is there, then we've
// hit an error.
func (parser *Parser) consume(tt token.Type, message string, args ...interface{}) {
	if parser.check(tt) {
		parser.advance()
		return
	}

	parser.newError(message, args...)
}

func (parser *Parser) checkNext(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.next().Type == tt
}

// isAtEnd checks if we've run out of tokens to parse.
func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

// peek returns the current token we have yet to consume.
func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.position]
}

// next returns the token ahead of the currently unconsumed token.
func (parser *Parser) next() token.Token {
	if !parser.isAtEnd() {
		return parser.tokens[parser.position+1]
	}

	return parser.tokens[parser.position]
}

// previous returns the most recently consumed token.
func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.position-1]
}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.peek().Type]; ok {
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

func (parser *Parser) newError(str string, args ...interface{}) {
	message := fmt.Sprintf(str, args...)

	parser.errors = append(parser.errors, message)
}
