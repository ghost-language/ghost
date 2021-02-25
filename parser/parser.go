package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

// Ghost uses a recursive decent parser. It starts from the top or outermost
// grammar rule (expression) and works its way down into the nested
// subexpressions before finally reaching the leaves of the syntax tree.

// A recursive decent parser is a literal translation of the grammar's rules
// straight into imperative code. Each rule becomes a function. The decent is
// described as "recursive" because when a grammar rule refers to itself --
// directly or indirectly -- that translates to a recursive function call.
// Recursive decent parsers are fast, robust, and can handle sophisticated
// error reporting while being easy to code and maintain.

// =============================================================================
// Precedence order

// statement              -> expressionStatement
// expressionStatement    -> expression
// expression             -> ternary ( equality "?" expression ":" expression )
// ternary                -> equality
// equality               -> comparison ( ( "!=" | "==" ) comparison )*
// comparison             -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
// term                   -> factor ( ( "-" | "+" ) factor )*
// factor                 -> unary ( ( "/" | "*" ) unary )*
// unary                  -> ( "!" | "-" ) unary
// primary                -> NUMBER | STRING | "true" | "false" | "null" |
//                        "(" expression ")"

// =============================================================================

// Parser - like the scanner - consumes a flat input sequence, only now we're
// reading tokens instead of characters. We store the list of tokens and use
// "current" to point to the next token eagerly waiting to be parsed.
type Parser struct {
	tokens  []token.Token
	current int
}

// New creates a new Parser instance.
func New(tokens []token.Token) Parser {
	return Parser{tokens, 0}
}

// Parse kicks off the parser.
func (parser *Parser) Parse() []ast.StatementNode {
	statements := make([]ast.StatementNode, 0)

	for !parser.isAtEnd() {
		statement := parser.statement()
		statements = append(statements, statement)
	}

	return statements
}

func (parser *Parser) statement() ast.StatementNode {
	if parser.match(token.PRINT) {
		return parser.printStatement()
	}

	return parser.expressionStatement()
}

func (parser *Parser) printStatement() ast.StatementNode {
	expression := parser.expression()

	return &ast.Print{Expression: expression}
}

func (parser *Parser) expressionStatement() ast.StatementNode {
	expression := parser.expression()

	return &ast.Expression{Expression: expression}
}

// expression starts the process of parsing expression grammar rules.

// Each method for parsing a grammar rule produces a syntax tree for that rule
// and returns it to the caller. When the body of the rule contains a non-
// terminal -- a reference to another rule -- we call that other rule's method.
func (parser *Parser) expression() ast.ExpressionNode {
	return parser.ternary()
}

func (parser *Parser) ternary() ast.ExpressionNode {
	condition := parser.equality()

	if parser.match("?") {
		thenExpression := parser.expression()

		parser.match(":")

		elseExpression := parser.expression()

		return &ast.Ternary{Condition: condition, Then: thenExpression, Else: elseExpression}
	}

	return condition
}

func (parser *Parser) equality() ast.ExpressionNode {
	expression := parser.comparison()

	for parser.match(token.BANGEQUAL, token.EQUALEQUAL) {
		operator := parser.previous()
		right := parser.comparison()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (parser *Parser) comparison() ast.ExpressionNode {
	expression := parser.term()

	if parser.match(token.GREATER, token.GREATEREQUAL, token.LESS, token.LESSEQUAL) {
		operator := parser.previous()
		right := parser.term()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (parser *Parser) term() ast.ExpressionNode {
	expression := parser.factor()

	for parser.match(token.MINUS, token.PLUS) {
		operator := parser.previous()
		right := parser.factor()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (parser *Parser) factor() ast.ExpressionNode {
	expression := parser.unary()

	for parser.match(token.SLASH, token.STAR) {
		operator := parser.previous()
		right := parser.unary()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (parser *Parser) unary() ast.ExpressionNode {
	if parser.match(token.BANG, token.MINUS) {
		operator := parser.previous()
		right := parser.unary()

		return &ast.Unary{Operator: operator, Right: right}
	}

	return parser.primary()
}

func (parser *Parser) primary() ast.ExpressionNode {
	if parser.match(token.FALSE) {
		return &ast.Boolean{Value: false}
	} else if parser.match(token.TRUE) {
		return &ast.Boolean{Value: true}
	} else if parser.match(token.NULL) {
		return &ast.Null{}
	} else if parser.match(token.NUMBER) {
		value, _ := decimal.NewFromString(parser.previous().Lexeme)
		return &ast.Number{Value: value}
	} else if parser.match(token.STRING) {
		return &ast.String{Value: parser.previous().Literal.(string)}
	} else if parser.match(token.LEFTPAREN) {
		expression := parser.expression()
		parser.consume(token.RIGHTPAREN, "Expected ')' after expression.")

		return &ast.Grouping{Expression: expression}
	}

	panic("Expected expression.")
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

func (parser *Parser) consume(tt token.Type, message string) token.Token {
	if parser.check(tt) {
		return parser.advance()
	}

	panic(fmt.Sprintf("%v: %v", parser.peek(), message))
}

// advance consumes the next token and pushes our current pointer ahead if we
// have not reached the end of our source code.
func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.current++
	}

	return parser.previous()
}

// check peeks ahead at the upcoming unconsumed token to verify it is of the
// type we passed through.
func (parser *Parser) check(tt token.Type) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == tt
}

// isAtEnd determines if we've reached the end of our source code.
func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

// peek looks at the current unconsumed token.
func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.current]
}

// previous returns the previous token.
func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}
