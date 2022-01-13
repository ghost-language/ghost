package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
)

// precedences contains a list of tokens mapped to their precedence level.
var precedences = map[token.Type]int{
	token.OR:           OR,
	token.AND:          AND,
	token.EQUALEQUAL:   EQUALS,
	token.BANGEQUAL:    EQUALS,
	token.IN:           EQUALS,
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
	token.LEFTBRACKET:  INDEX,
	token.DOT:          INDEX,
	token.PLUSEQUAL:    SUM,
	token.MINUSEQUAL:   SUM,
	token.STAREQUAL:    PRODUCT,
	token.SLASHEQUAL:   PRODUCT,
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
	scanner *scanner.Scanner
	errors  []string

	previousToken token.Token
	currentToken  token.Token
	nextToken     token.Token

	previousIndex    *ast.Index
	previousProperty *ast.Property

	prefixParserFns  map[token.Type]prefixParserFn
	infixParserFns   map[token.Type]infixParserFn
	postfixParserFns map[token.Type]postfixParserFn
}

// New creates a new parser instance.
func New(scanner *scanner.Scanner) *Parser {
	parser := &Parser{
		scanner:          scanner,
		errors:           []string{},
		prefixParserFns:  make(map[token.Type]prefixParserFn),
		infixParserFns:   make(map[token.Type]infixParserFn),
		postfixParserFns: make(map[token.Type]postfixParserFn),
	}

	// Register all of our prefix parse functions
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
	parser.registerPrefix(token.LEFTBRACKET, parser.listLiteral)
	parser.registerPrefix(token.LEFTBRACE, parser.mapLiteral)
	parser.registerPrefix(token.WHILE, parser.whileExpression)
	parser.registerPrefix(token.FOR, parser.forExpression)
	parser.registerPrefix(token.CLASS, parser.classStatement)
	parser.registerPrefix(token.THIS, parser.thisExpression)

	// Register all of our infix parse functions
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
	parser.registerInfix(token.LEFTBRACKET, parser.indexExpression)
	parser.registerInfix(token.DOT, parser.dotExpression)
	parser.registerInfix(token.AND, parser.infixExpression)
	parser.registerInfix(token.OR, parser.infixExpression)
	parser.registerInfix(token.PLUSEQUAL, parser.compoundExpression)
	parser.registerInfix(token.MINUSEQUAL, parser.compoundExpression)
	parser.registerInfix(token.STAREQUAL, parser.compoundExpression)
	parser.registerInfix(token.SLASHEQUAL, parser.compoundExpression)

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

		program.Statements = append(program.Statements, statement)

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
	parser.previousToken = parser.currentToken
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.scanner.ScanToken()
}

// // isAtEnd checks if we've run out of tokens to parse.
func (parser *Parser) isAtEnd() bool {
	return parser.currentTokenIs(token.EOF)
}

func (parser *Parser) nextError(tt token.Type) {
	message := fmt.Sprintf(
		"%d:%d: syntax error: expected next token to be %s, got: %s instead", parser.nextToken.Line, parser.nextToken.Column, tt, parser.nextToken.Type,
	)

	parser.errors = append(parser.errors, message)
}

func (parser *Parser) currentTokenIs(tt token.Type) bool {
	return parser.currentToken.Type == tt
}

func (parser *Parser) nextTokenIs(tt token.Type) bool {
	return parser.nextToken.Type == tt
}

func (parser *Parser) expectNextTokenIs(tt token.Type) bool {
	if parser.nextTokenIs(tt) {
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
