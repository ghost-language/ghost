package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/ghost"
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

// statement              -> expressionStatement | printStatement |
//                        blockStatement
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
		statement, _ := parser.statement()

		statements = append(statements, statement)
	}

	return statements
}

func (parser *Parser) statement() (ast.StatementNode, error) {
	if parser.match(token.PRINT) {
		return parser.printStatement()
	} else if parser.match(token.LEFTBRACE) {
		var err error

		if statements, err := parser.block(); err == nil {
			return &ast.Block{Statements: statements}, nil
		}

		return nil, err
	}

	return parser.expressionStatement()
}

func (parser *Parser) printStatement() (ast.StatementNode, error) {
	expression, err := parser.expression()

	if err != nil {
		return nil, err
	}

	return &ast.Print{Expression: expression}, nil
}

func (parser *Parser) block() ([]ast.StatementNode, error) {
	statements := make([]ast.StatementNode, 0)

	for !parser.check(token.RIGHTBRACE) && !parser.isAtEnd() {
		statement, err := parser.statement()

		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	parser.consume(token.RIGHTBRACE, "Expected closing bracket '}' after block.")

	return statements, nil
}

func (parser *Parser) expressionStatement() (ast.StatementNode, error) {
	expression, err := parser.expression()

	if err != nil {
		return nil, err
	}

	return &ast.Expression{Expression: expression}, nil
}

// expression starts the process of parsing expression grammar rules.

// Each method for parsing a grammar rule produces a syntax tree for that rule
// and returns it to the caller. When the body of the rule contains a non-
// terminal -- a reference to another rule -- we call that other rule's method.
func (parser *Parser) expression() (ast.ExpressionNode, error) {
	return parser.assign()
}

func (parser *Parser) assign() (ast.ExpressionNode, error) {
	if parser.check(token.IDENTIFIER) && parser.next().Type == token.ASSIGN {
		parser.match(token.IDENTIFIER)
		name := parser.previous()

		parser.match(token.ASSIGN)

		val, err := parser.expression()

		if err != nil {
			return nil, err
		}

		return &ast.Assign{Name: name, Value: val}, nil
	}

	expression, err := parser.ternary()

	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (parser *Parser) ternary() (ast.ExpressionNode, error) {
	condition, err := parser.equality()

	if err != nil {
		return nil, err
	}

	if parser.match("?") {
		thenExpression, err := parser.expression()

		if err != nil {
			return nil, err
		}

		parser.match(":")

		elseExpression, err := parser.expression()

		if err != nil {
			return nil, err
		}

		return &ast.Ternary{Condition: condition, Then: thenExpression, Else: elseExpression}, nil
	}

	return condition, nil
}

func (parser *Parser) equality() (ast.ExpressionNode, error) {
	expression, err := parser.comparison()

	if err != nil {
		return nil, err
	}

	for parser.match(token.BANGEQUAL, token.EQUALEQUAL) {
		operator := parser.previous()
		right, err := parser.comparison()

		if err != nil {
			return nil, err
		}

		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (parser *Parser) comparison() (ast.ExpressionNode, error) {
	expression, err := parser.term()

	if err != nil {
		return nil, err
	}

	if parser.match(token.GREATER, token.GREATEREQUAL, token.LESS, token.LESSEQUAL) {
		operator := parser.previous()
		right, err := parser.term()

		if err != nil {
			return nil, err
		}

		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (parser *Parser) term() (ast.ExpressionNode, error) {
	expression, err := parser.factor()

	if err != nil {
		return nil, err
	}

	for parser.match(token.MINUS, token.PLUS) {
		operator := parser.previous()
		right, err := parser.factor()

		if err != nil {
			return nil, err
		}

		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (parser *Parser) factor() (ast.ExpressionNode, error) {
	expression, err := parser.unary()

	if err != nil {
		return nil, err
	}

	for parser.match(token.SLASH, token.STAR) {
		operator := parser.previous()
		right, err := parser.unary()

		if err != nil {
			return nil, err
		}

		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (parser *Parser) unary() (ast.ExpressionNode, error) {
	if parser.match(token.BANG, token.MINUS) {
		operator := parser.previous()
		right, err := parser.unary()

		if err != nil {
			return nil, err
		}

		return &ast.Unary{Operator: operator, Right: right}, nil
	}

	return parser.primary()
}

func (parser *Parser) primary() (ast.ExpressionNode, error) {
	if parser.match(token.FALSE) {
		return &ast.Boolean{Value: false}, nil
	} else if parser.match(token.TRUE) {
		return &ast.Boolean{Value: true}, nil
	} else if parser.match(token.NULL) {
		return &ast.Null{}, nil
	} else if parser.match(token.NUMBER) {
		value, _ := decimal.NewFromString(parser.previous().Lexeme)
		return &ast.Number{Value: value}, nil
	} else if parser.match(token.STRING) {
		return &ast.String{Value: parser.previous().Literal.(string)}, nil
	} else if parser.match(token.LEFTPAREN) {
		var expression ast.ExpressionNode
		var err error

		if parser.check(token.RIGHTPAREN) {
			expression = &ast.Null{}
		} else {
			expression, err = parser.expression()

			if err != nil {
				return nil, err
			}
		}

		_, err = parser.consume(token.RIGHTPAREN, "Expected ')' after expression.")

		if err != nil {
			return nil, err
		}

		return &ast.Grouping{Expression: expression}, nil
	} else if parser.match(token.IDENTIFIER) {
		return &ast.Variable{Name: parser.previous()}, nil
	}

	return nil, ghost.ParseError(parser.peek(), fmt.Sprintf("Expected expression, got=%v", parser.peek().Type))
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

func (parser *Parser) consume(tt token.Type, message string) (token.Token, error) {
	if parser.check(tt) {
		return parser.advance(), nil
	}

	return parser.previous(), ghost.ParseError(parser.peek(), message)
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

func (parser *Parser) next() token.Token {
	return parser.tokens[parser.current+1]
}

// previous returns the previous token.
func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}
