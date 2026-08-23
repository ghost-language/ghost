package parser

import (
	"sort"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
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
	token.DOTDOT:       RANGE,
	token.QUESTION:     TERNARY,

	// `++`/`--` bind at the same precedence as `.` and `[]` (not tighter), so
	// that when they trail a property or index chain, the chain's own
	// recursive parse (e.g. dotExpression parsing the property name at INDEX
	// precedence) stops before consuming the operator itself, leaving it for
	// the outer loop to attach to the whole chain rather than just its last
	// segment.
	token.PLUSPLUS:   INDEX,
	token.MINUSMINUS: INDEX,
}

// The following list of constants define the available precedence levels.
const (
	_ int = iota
	LOWEST
	OR
	AND
	TERNARY
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
	prefixParserFn func() ast.ExpressionNode
	infixParserFn  func(ast.ExpressionNode) ast.ExpressionNode
)

// Parser holds a slice of tokens, its position, and errors
// as well as the prefix and infix parse functions.
type Parser struct {
	scanner *scanner.Scanner
	errors  []*fault.Fault

	previousToken token.Token
	currentToken  token.Token
	nextToken     token.Token

	previousIndex    *ast.Index
	previousProperty *ast.Property

	prefixParserFns map[token.Type]prefixParserFn
	infixParserFns  map[token.Type]infixParserFn

	inTernaryExpression bool
}

// New creates a new parser instance.
func New(scanner *scanner.Scanner) *Parser {
	parser := &Parser{
		scanner:         scanner,
		errors:          []*fault.Fault{},
		prefixParserFns: make(map[token.Type]prefixParserFn),
		infixParserFns:  make(map[token.Type]infixParserFn),
	}

	// Register all of our prefix parse functions
	parser.registerPrefix(token.IDENTIFIER, parser.identifierLiteral)
	parser.registerPrefix(token.NUMBER, parser.numberLiteral)
	parser.registerPrefix(token.NULL, parser.nullLiteral)
	parser.registerPrefix(token.TRUE, parser.booleanLiteral)
	parser.registerPrefix(token.FALSE, parser.booleanLiteral)
	parser.registerPrefix(token.STRING, parser.stringLiteral)
	parser.registerPrefix(token.TEMPLATESTRING, parser.templateLiteral)
	parser.registerPrefix(token.TEMPLATESTRINGEND, parser.templateLiteral)
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
	parser.registerPrefix(token.TRAIT, parser.traitStatement)
	parser.registerPrefix(token.USE, parser.useExpression)
	parser.registerPrefix(token.THIS, parser.thisExpression)
	parser.registerPrefix(token.SUPER, parser.superExpression)
	parser.registerPrefix(token.NEW, parser.newExpression)
	parser.registerPrefix(token.IMPORT, parser.importStatement)
	parser.registerPrefix(token.SWITCH, parser.switchStatement)
	parser.registerPrefix(token.BREAK, parser.breakStatement)
	parser.registerPrefix(token.CONTINUE, parser.continueStatement)

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
	parser.registerInfix(token.DOTDOT, parser.infixExpression)
	parser.registerInfix(token.PLUSEQUAL, parser.compoundExpression)
	parser.registerInfix(token.MINUSEQUAL, parser.compoundExpression)
	parser.registerInfix(token.STAREQUAL, parser.compoundExpression)
	parser.registerInfix(token.SLASHEQUAL, parser.compoundExpression)
	parser.registerInfix(token.QUESTION, parser.ternaryExpression)
	parser.registerInfix(token.PLUSPLUS, parser.postfixExpression)
	parser.registerInfix(token.MINUSMINUS, parser.postfixExpression)

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
		reported := len(parser.errors)
		statement := parser.statement()

		program.Statements = append(program.Statements, statement)

		// A statement that failed leaves the parser somewhere in the middle of
		// it, where every following token is read out of context. Skipping to
		// the next statement boundary keeps one mistake from being reported as
		// a dozen.
		if len(parser.errors) > reported {
			parser.synchronize()

			continue
		}

		parser.readToken()
	}

	return program
}

// Errors returns everything wrong with the source, in the order a reader would
// come across it. Lexical and grammatical problems are one list rather than
// two: which pass noticed a problem is Ghost's business, not the reader's.
func (parser *Parser) Errors() []*fault.Fault {
	raised := make([]*fault.Fault, 0, len(parser.errors)+len(parser.scanner.Faults()))

	raised = append(raised, parser.scanner.Faults()...)
	raised = append(raised, parser.errors...)

	sort.SliceStable(raised, func(first int, second int) bool {
		left := raised[first].Position
		right := raised[second].Position

		if left.Line != right.Line {
			return left.Line < right.Line
		}

		return left.Column < right.Column
	})

	return raised
}

// synchronize discards tokens until the parser is at something that can begin a
// statement again. Recovering this way is what lets a single run report every
// syntax error in a file instead of only the first.
func (parser *Parser) synchronize() {
	for !parser.isAtEnd() {
		if parser.currentTokenIs(token.SEMICOLON) {
			parser.readToken()

			return
		}

		switch parser.nextToken.Type {
		case token.CLASS, token.FUNCTION, token.FOR, token.IF, token.RETURN,
			token.SWITCH, token.TRAIT, token.WHILE, token.IMPORT, token.USE:
			parser.readToken()

			return
		}

		parser.readToken()
	}
}

// report records a syntax error at a token.
func (parser *Parser) report(tok token.Token, format string, arguments ...interface{}) *fault.Fault {
	raised := fault.At(fault.Syntax, tok, format, arguments...)

	parser.errors = append(parser.errors, raised)

	return raised
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
	parser.report(parser.nextToken, "expected %s, found %s", tt.Describe(), parser.nextToken.Describe())
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

// prefixError records a syntax error for a token that cannot begin an
// expression. Without it an unparseable token yields a nil expression that only
// surfaces later as a crash in the evaluator.
func (parser *Parser) prefixError() {
	raised := parser.report(parser.currentToken, "%s cannot start an expression", capitalize(parser.currentToken.Describe()))

	// A closing bracket here almost always means the matching opener was never
	// closed, and saying so is more use than naming the token that tripped.
	switch parser.currentToken.Type {
	case token.RIGHTPAREN, token.RIGHTBRACKET, token.RIGHTBRACE:
		raised.WithHelp("check for a missing opening `%s` or an extra `%s`", opener(parser.currentToken.Type), parser.currentToken.Lexeme)
	}
}

// opener names the bracket that should have preceded a closing one.
func opener(closing token.Type) string {
	switch closing {
	case token.RIGHTPAREN:
		return "("
	case token.RIGHTBRACKET:
		return "["
	}

	return "{"
}

// capitalize starts a sentence with an upper-case letter without disturbing the
// backticks a description may begin with.
func capitalize(text string) string {
	for index, character := range text {
		if character >= 'a' && character <= 'z' {
			return text[:index] + string(character-32) + text[index+1:]
		}

		if character == '`' {
			break
		}
	}

	return text
}
