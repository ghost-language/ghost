package parser

import (
	"fmt"
	"log"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/errors"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
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

// parse                  -> declaration
// declaration            -> classDeclaration | letDeclaration | functionDeclaration | statement
// classDeclaration       -> "class" IDENTIFIER "{" function* "}"
// functionDeclaration    -> "function" IDENTIFIER "(" parameters? ")" block
// letDeclaration         -> "let" IDENTIFIER ( "=" expression )
// statement              -> expressionStatement | ifStatement | whileStatement |
//                        printStatement | blockStatement
// expressionStatement    -> expression
// expression             -> or
// or                     -> and ( "or" and )
// and                    -> ternary ( "and" ternary )
// ternary                -> equality ( equality "?" expression ":" expression )
// equality               -> comparison ( ( "!=" | "==" ) comparison )*
// comparison             -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
// term                   -> factor ( ( "-" | "+" ) factor )*
// factor                 -> unary ( ( "/" | "*" ) unary )*
// unary                  -> ( "!" | "-" ) unary | call
// call                   -> primary ( "(" arguments? ")" )
// arguments              -> expression ( "," expression )
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
		statement, _ := parser.declaration()

		statements = append(statements, statement)
	}

	return statements
}

func (parser *Parser) declaration() (ast.StatementNode, error) {
	var statement ast.StatementNode
	var err error

	if parser.match(token.CLASS) {
		statement, err = parser.classDeclaration()
	} else if parser.match(token.FUNCTION) {
		statement, err = parser.functionDeclaration("function")
	} else {
		statement, err = parser.statement()
	}

	return statement, err
}

func (parser *Parser) classDeclaration() (ast.StatementNode, error) {
	name, err := parser.consume(token.IDENTIFIER, "Expected class name.")

	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.LEFTBRACE, "Expected '{' before class body.")

	if err != nil {
		return nil, err
	}

	methods := make([]*ast.Function, 0)

	for !parser.check(token.RIGHTBRACE) && !parser.isAtEnd() {
		function, err := parser.functionDeclaration("method")

		if err != nil {
			return nil, err
		}

		methods = append(methods, function)
	}

	_, err = parser.consume(token.RIGHTBRACE, "Expected '}' after class body.")

	if err != nil {
		return nil, err
	}

	return &ast.Class{Name: name, Methods: methods}, nil
}

func (parser *Parser) letDeclaration() (ast.StatementNode, error) {
	name, err := parser.consume(token.IDENTIFIER, "Expected variable name.")

	if err != nil {
		return nil, err
	}

	var initializer ast.ExpressionNode

	if parser.match(token.EQUAL) {
		initializer, err = parser.expression()

		if err != nil {
			return nil, err
		}
	}

	// consume the semicolon and continue on
	parser.match(token.SEMICOLON)

	return &ast.Declaration{Name: name, Initializer: initializer}, nil
}

func (parser *Parser) functionDeclaration(kind string) (*ast.Function, error) {
	name, err := parser.consume(token.IDENTIFIER, "Expected " + kind + " name.")

	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.LEFTPAREN, "Expected '(' after " + kind + " name.")

	if err != nil {
		return nil, err
	}

	parameters := make([]token.Token, 0)

	if !parser.check(token.RIGHTPAREN) {
		for {
			param, err := parser.consume(token.IDENTIFIER, "Expected parameter name.")

			if err != nil {
				return nil, err
			}

			parameters = append(parameters, param)

			if !parser.match(token.COMMA) {
				break
			}
		}
	}

	_, err = parser.consume(token.RIGHTPAREN, "Expected ')' after parameters.")

	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.LEFTBRACE, "Expected '{' before " + kind + " body.")

	if err != nil {
		return nil, err
	}

	body, err := parser.block()

	if err != nil {
		return nil, err
	}

	log.Fatal(fmt.Sprintf("%s : %s", name.Lexeme, body))

	return nil, nil
	// return &ast.Function{Name: name, Params: parameters, Body: body}, nil
}

func (parser *Parser) statement() (ast.StatementNode, error) {
	if parser.match(token.FOR) {
		return parser.forStatement()
	} else if parser.match(token.IF) {
		return parser.ifStatement()
	} else if parser.match(token.WHILE) {
		return parser.whileStatement()
	} else if parser.match(token.PRINT) {
		return parser.printStatement()
	} else if parser.match(token.LEFTBRACE) {
		statements, err := parser.block()

		if err != nil {
			return nil, err
		}

		return &ast.Block{Statements: statements}, nil
	}

	return parser.expressionStatement()
}

func (parser *Parser) forStatement() (ast.StatementNode, error) {
	_, err := parser.consume(token.LEFTPAREN, "Expected '(' after 'for'.")

	if err != nil {
		return nil, err
	}

	// Initializer
	var initializer ast.StatementNode

	if parser.match(token.SEMICOLON) {
		initializer = nil
	} else if parser.match(token.LET) {
		initializer, err = parser.letDeclaration()

		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = parser.expressionStatement()

		if err != nil {
			return nil, err
		}
	}

	// Condition
	var condition ast.ExpressionNode

	if !parser.check(token.SEMICOLON) {
		condition, err = parser.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(token.SEMICOLON, "Expect ';' afer loop condition.")

	if err != nil {
		return nil, err
	}

	// Increment
	var increment ast.ExpressionNode

	if !parser.check(token.RIGHTPAREN) {
		increment, err = parser.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(token.RIGHTPAREN, "Expected ')' after for clause.")

	if err != nil {
		return nil, err
	}

	// Body
	body, err := parser.statement()

	if err != nil {
		return nil, err
	}

	// Desugaring to while loop 🦾
	if increment != nil {
		statements := make([]ast.StatementNode, 0)
		statements = append(statements, body)
		statements = append(statements, increment)

		body = &ast.Block{Statements: statements}
	}

	if condition == nil {
		condition = value.TRUE
	}

	body = &ast.While{Condition: condition, Body: body}

	if initializer != nil {
		statements := make([]ast.StatementNode, 0)
		statements = append(statements, initializer)
		statements = append(statements, body)

		body = &ast.Block{Statements: statements}
	}

	return body, nil
}

func (parser *Parser) ifStatement() (ast.StatementNode, error) {
	var err error

	if _, err := parser.consume(token.LEFTPAREN, "Expected '(' after 'if'."); err != nil {
		return nil, err
	}

	if condition, err := parser.expression(); err == nil {
		if _, err := parser.consume(token.RIGHTPAREN, "Expected ')' after 'if' condition."); err == nil {
			if thenBranch, err := parser.statement(); err == nil {
				if parser.match(token.ELSE) {
					if elseBranch, err := parser.statement(); err == nil {
						return &ast.If{Condition: condition, Then: thenBranch, Else: elseBranch}, nil
					}
				} else {
					return &ast.If{Condition: condition, Then: thenBranch}, nil
				}
			}
		}
	}

	return nil, err
}

func (parser *Parser) whileStatement() (ast.StatementNode, error) {
	_, err := parser.consume(token.LEFTPAREN, "Expected '(' after 'while'.")

	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()

	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.RIGHTPAREN, "Expected ')' after 'while' condition.")

	if err != nil {
		return nil, err
	}

	body, err := parser.statement()

	if err != nil {
		return nil, err
	}

	return &ast.While{Condition: condition, Body: body}, nil
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
		statement, err := parser.declaration()

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

	expression, err := parser.or()

	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (parser *Parser) or() (ast.ExpressionNode, error) {
	expression, err := parser.and()

	if err != nil {
		return nil, err
	}

	for parser.match(token.OR) {
		operator := parser.previous()
		right, err := parser.and()

		if err != nil {
			return nil, err
		}

		expression = &ast.Logical{Left: expression, Operator: operator, Right: right}
	}

	return expression, err
}

func (parser *Parser) and() (ast.ExpressionNode, error) {
	expression, err := parser.ternary()

	if err != nil {
		return nil, err
	}

	for parser.match(token.AND) {
		operator := parser.previous()
		right, err := parser.ternary()

		if err != nil {
			return nil, err
		}

		expression = &ast.Logical{Left: expression, Operator: operator, Right: right}
	}

	return expression, err
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

	return parser.call()
}

func (parser *Parser) call() (ast.ExpressionNode, error) {
	expression, err := parser.primary()

	if err != nil {
		return nil, err
	}

	for {
		if parser.match(token.LEFTPAREN) {
			expression, err = parser.finishCall(expression)

			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expression, err
}

func (parser *Parser) finishCall(callee ast.ExpressionNode) (ast.ExpressionNode, error) {
	arguments := make([]ast.ExpressionNode, 0)

	if !parser.check(token.RIGHTPAREN) {
		for {
			argument, err := parser.assign()

			if err != nil {
				return nil, err
			}

			if len(arguments) >= 255 {
				return nil, errors.ParseError(parser.peek(), "Cannot have more than 255 arguments.")
			}

			arguments = append(arguments, argument)

			if !parser.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := parser.consume(token.RIGHTPAREN, "Expected ')' after arguments.")

	if err != nil {
		return nil, err
	}

	return &ast.Call{Callee: callee, Paren: paren, Arguments: arguments}, nil
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
		return &ast.Identifier{Name: parser.previous()}, nil
	}

	return nil, errors.ParseError(parser.peek(), fmt.Sprintf("Expected expression, got=%v", parser.peek().Type))
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

	return parser.previous(), errors.ParseError(parser.peek(), message)
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
