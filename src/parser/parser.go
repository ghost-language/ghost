package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// precedences contains a list of tokens mapped to their precedence level.
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

// The following list of constants define the available precedence levels.
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

// Parser holds a slice of tokens, its position, and errors
// as well as the prefix, infix, and postfix parse functions.
type Parser struct {
	tokens   []token.Token
	position int
	errors   []string

	previousToken token.Token
	currentToken  token.Token
	nextToken     token.Token

	prefixParserFns  map[token.Type]prefixParserFn
	infixParserFns   map[token.Type]infixParserFn
	postfixParserFns map[token.Type]postfixParserFn
}

// New creates a new parser instance.
func New(tokens []token.Token) *Parser {
	parser := &Parser{
		tokens:           tokens,
		position:         0,
		errors:           []string{},
		prefixParserFns:  make(map[token.Type]prefixParserFn),
		infixParserFns:   make(map[token.Type]infixParserFn),
		postfixParserFns: make(map[token.Type]postfixParserFn),
	}

	// Register all of our prefix parse functions
	parser.registerPrefix(token.IDENTIFIER, parser.identifierLiteral)
	// parser.registerPrefix(token.NUMBER, parser.numberLiteral)
	// parser.registerPrefix(token.NULL, parser.nullLiteral)
	parser.registerPrefix(token.TRUE, parser.booleanLiteral)
	parser.registerPrefix(token.FALSE, parser.booleanLiteral)
	// parser.registerPrefix(token.STRING, parser.stringLiteral)
	// parser.registerPrefix(token.BANG, parser.prefixExpression)
	// parser.registerPrefix(token.MINUS, parser.prefixExpression)
	// parser.registerPrefix(token.IF, parser.ifExpression)
	// parser.registerPrefix(token.LEFTPAREN, parser.groupExpression)
	// parser.registerPrefix(token.FUNCTION, parser.functionStatement)

	// Register all of our infix parse functions
	// parser.registerInfix(token.PLUS, parser.infixExpression)
	// parser.registerInfix(token.MINUS, parser.infixExpression)
	// parser.registerInfix(token.SLASH, parser.infixExpression)
	// parser.registerInfix(token.STAR, parser.infixExpression)
	// parser.registerInfix(token.PERCENT, parser.infixExpression)
	// parser.registerInfix(token.EQUALEQUAL, parser.infixExpression)
	// parser.registerInfix(token.BANGEQUAL, parser.infixExpression)
	// parser.registerInfix(token.GREATER, parser.infixExpression)
	// parser.registerInfix(token.GREATEREQUAL, parser.infixExpression)
	// parser.registerInfix(token.LESS, parser.infixExpression)
	// parser.registerInfix(token.LESSEQUAL, parser.infixExpression)
	// parser.registerInfix(token.LEFTPAREN, parser.callExpression)

	// Read the first two tokens, so currentToken and nextToken are both set.
	parser.readToken()
	parser.readToken()

	return parser
}

// registerPrefix registers a new prefix parse function.
func (parser *Parser) registerPrefix(tokenType token.Type, fn prefixParserFn) {
	parser.prefixParserFns[tokenType] = fn
}

// registerInfix registers a new infix parse function.
func (parser *Parser) registerInfix(tokenType token.Type, fn infixParserFn) {
	parser.infixParserFns[tokenType] = fn
}

// Parse parses tokens and creates an AST. It returns the Program node,
// which holds a slice of Statements (and in turn, the rest of the tree).
func (parser *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.StatementNode{}

	for !parser.isAtEnd() {
		statement := parser.statement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.readToken()
	}

	return program
}

// Errors returns the slice of errors contained within the parser instance.
func (parser *Parser) Errors() []string {
	return parser.errors
}

// =============================================================================
// Helper methods

// readToken advances the parser through the list of tokens, setting the
// previous, current, and next token values for consumption.
func (parser *Parser) readToken() {
	if !parser.isAtEnd() {
		parser.previousToken = parser.currentToken
		parser.currentToken = parser.nextToken

		if parser.position >= 0 && parser.position < len(parser.tokens) {
			parser.nextToken = parser.tokens[parser.position]
		}

		parser.position++
	}
}

// // isAtEnd checks if we've run out of tokens to parse.
func (parser *Parser) isAtEnd() bool {
	return parser.currentTokenTypeIs(token.EOF)
}

func (parser *Parser) nextError(tt token.Type) {
	message := fmt.Sprintf(
		"Line: %d: Expected next token to be %s, got: %s instead", parser.currentToken.Line, tt, parser.nextToken.Type,
	)

	parser.errors = append(parser.errors, message)
}

func (parser *Parser) currentTokenTypeIs(tt token.Type) bool {
	return parser.currentToken.Type == tt
}

func (parser *Parser) nextTokenTypeIs(tt token.Type) bool {
	return parser.nextToken.Type == tt
}

func (parser *Parser) expectNextType(tt token.Type) bool {
	if parser.nextTokenTypeIs(tt) {
		parser.readToken()
		return true
	}

	parser.nextError(tt)
	return false
}

func (parser *Parser) nextTokenPrecedence() int {
	if precedence, ok := precedences[parser.nextToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) currentTokenPrecedence() int {
	if precedence, ok := precedences[parser.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

// // match checks to see if the current token has any of the given types. If so,
// // it consumes the token and returns true. Otherwise, it returns false and
// // leaves the current token alone.
// func (parser *Parser) match(tt ...token.Type) bool {
// 	for _, t := range tt {
// 		if parser.check(t) {
// 			parser.advance()
// 			return true
// 		}
// 	}

// 	return false
// }

// // advance consumes the current token and returns it, similar to how our
// // scanner's corresponding method crawls through characters.
// func (parser *Parser) advance() token.Token {
// 	if !parser.isAtEnd() {
// 		parser.position++
// 	}

// 	return parser.previous()
// }

// // check returns true if the current token is of the given type. Unlike match(),
// // it never consumes the token, it only looks at it.
// func (parser *Parser) check(tt token.Type) bool {
// 	if parser.isAtEnd() {
// 		return false
// 	}

// 	return parser.peek().Type == tt
// }

// // consume checks to see if the next token is of the expected typ. If so, it
// // consumes the token and carries on. If some other token is there, then we've
// // hit an error.
// func (parser *Parser) consume(tt token.Type, message string, args ...interface{}) {
// 	if parser.check(tt) {
// 		parser.advance()
// 		return
// 	}

// 	parser.newError(message, args...)
// }

// func (parser *Parser) checkNext(tt token.Type) bool {
// 	if parser.isAtEnd() {
// 		return false
// 	}

// 	return parser.next().Type == tt
// }

// // peek returns the current token we have yet to consume.
// func (parser *Parser) peek() token.Token {
// 	return parser.tokens[parser.position]
// }

// // next returns the token ahead of the currently unconsumed token.
// func (parser *Parser) next() token.Token {
// 	if !parser.isAtEnd() {
// 		return parser.tokens[parser.position+1]
// 	}

// 	return parser.tokens[parser.position]
// }

// // previous returns the most recently consumed token.
// func (parser *Parser) previous() token.Token {
// 	return parser.tokens[parser.position-1]
// }

// func (parser *Parser) peekPrecedence() int {
// 	if precedence, ok := precedences[parser.peek().Type]; ok {
// 		return precedence
// 	}

// 	return LOWEST
// }

// func (parser *Parser) nextPrecedence() int {
// 	if parser, ok := precedences[parser.next().Type]; ok {
// 		return parser
// 	}

// 	return LOWEST
// }

// func (parser *Parser) newError(str string, args ...interface{}) {
// 	message := fmt.Sprintf(str, args...)

// 	parser.errors = append(parser.errors, message)
// }
