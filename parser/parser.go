package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

var precedences = map[token.TokenType]int{
	token.OR:             OR,
	token.AND:            AND,
	token.EQ:             EQUALS,
	token.NOTEQ:          EQUALS,
	token.IN:             EQUALS,
	token.COMMA:          EQUALS,
	token.LT:             LESSGREATER,
	token.GT:             LESSGREATER,
	token.LTE:            LESSGREATER,
	token.GTE:            LESSGREATER,
	token.PLUS:           SUM,
	token.MINUS:          SUM,
	token.PLUSASSIGN:     SUM,
	token.MINUSASSIGN:    SUM,
	token.SLASH:          PRODUCT,
	token.ASTERISK:       PRODUCT,
	token.SLASHASSIGN:    PRODUCT,
	token.ASTERISKASSIGN: PRODUCT,
	token.PERCENT:        MODULO,
	token.LPAREN:         CALL,
	token.LBRACKET:       INDEX,
	token.DOT:            INDEX,
	token.RANGE:          RANGE,
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
	prefixParserFn func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func() ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	previousToken token.Token
	currentToken  token.Token
	peekToken     token.Token

	prefixParseFns  map[token.TokenType]prefixParserFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn

	previousIndexExpression    *ast.IndexExpression
	previousPropertyExpression *ast.PropertyExpression
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParserFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifierLiteral)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseListLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.IMPORT, p.parseImportExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.IN, p.parseInfixExpression)
	p.registerInfix(token.PLUSASSIGN, p.parseInfixExpression)
	p.registerInfix(token.MINUSASSIGN, p.parseInfixExpression)
	p.registerInfix(token.ASTERISKASSIGN, p.parseInfixExpression)
	p.registerInfix(token.SLASHASSIGN, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseDotExpression)
	p.registerInfix(token.RANGE, p.parseInfixExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.PLUSPLUS, p.parsePostfixExpression)
	p.registerPostfix(token.MINUSMINUS, p.parsePostfixExpression)

	// Read two tokens, so currentToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParserFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.previousToken = p.currentToken
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if p.currentToken.Type == token.RETURN {
		return p.parseReturnStatement()
	}

	statement := p.parseAssignStatement()

	if statement != nil {
		return statement
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	postfix := p.postfixParseFns[p.currentToken.Type]

	if postfix != nil {
		return postfix()
	}

	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

// currentTokenIs determines if the current token is of the type specified.
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs determines if the next token is of the type specified.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek determines if the next token is of the type specified.
// If not, an error will be thrown.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)

	p.errors = append(p.errors, message)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, message)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currentPrecendence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

// ----------------------------------------------------------------------------
// Literals

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parseIdentifierLiteral() ast.Expression {
	return &ast.IdentifierLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

// Anonymous function:  function() { ... }
// Named function:      function test() { ... }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.peekTokenIs(token.LPAREN) {
		p.nextToken()

		literal.Name = p.currentToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Defaults, literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseFunctionParameters() (map[string]ast.Expression, []*ast.IdentifierLiteral) {
	defaults := make(map[string]ast.Expression)
	identifiers := []*ast.IdentifierLiteral{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()

		return defaults, identifiers
	}

	p.nextToken()

	for !p.currentTokenIs(token.RPAREN) {
		identifier := &ast.IdentifierLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, identifier)

		p.nextToken()

		if p.currentTokenIs(token.ASSIGN) {
			p.nextToken()
			defaults[identifier.Value] = p.parseExpressionStatement().Expression
			p.nextToken()
		}

		if p.currentTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return defaults, identifiers
}

func (p *Parser) parseListLiteral() ast.Expression {
	list := &ast.ListLiteral{Token: p.currentToken}
	list.Elements = p.parseExpressionList(token.RBRACKET)

	return list
}

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLiteral := &ast.MapLiteral{Token: p.currentToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)
		mapLiteral.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return mapLiteral
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	numberLiteral := &ast.NumberLiteral{Token: p.currentToken}

	value, err := decimal.NewFromString(p.currentToken.Literal)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as number", p.currentToken.Literal)
		p.errors = append(p.errors, message)

		return nil
	}

	numberLiteral.Value = value

	return numberLiteral
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

// ----------------------------------------------------------------------------
// Expressions

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression() ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.previousToken,
		Operator: p.currentToken.Literal,
	}
}

func (p *Parser) parseCallExpression(callable ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Callable: callable}
	expression.Arguments = p.parseExpressionList(token.RPAREN)

	return expression
}

func (p *Parser) parseDotExpression(object ast.Expression) ast.Expression {
	tok := p.currentToken

	precedence := p.currentPrecendence()
	p.nextToken()

	if p.peekTokenIs(token.LPAREN) {
		// method
		exp := &ast.MethodExpression{Token: tok, Object: object}
		exp.Method = p.parseExpression(precedence)
		p.nextToken()
		exp.Arguments = p.parseExpressionList(token.RPAREN)

		return exp
	}

	// property
	exp := &ast.PropertyExpression{Token: tok, Object: object}
	exp.Property = p.parseIdentifierLiteral()

	p.previousPropertyExpression = exp
	p.previousIndexExpression = nil

	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	if !p.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return p.parseForInExpression(expression)
	}

	expression.Identifier = p.currentToken.Literal
	expression.Initializer = p.parseAssignStatement()

	if expression.Initializer == nil {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if expression.Condition == nil {
		return nil
	}

	p.nextToken()
	p.nextToken()

	expression.Increment = p.parseAssignStatement()

	if expression.Increment == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Block = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForInExpression(parentExpression *ast.ForExpression) ast.Expression {
	expression := &ast.ForInExpression{Token: parentExpression.Token}

	if !p.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	value := p.currentToken.Literal
	var key string
	p.nextToken()

	if p.currentTokenIs(token.COMMA) {
		p.nextToken()

		if !p.currentTokenIs(token.IDENTIFIER) {
			return nil
		}

		key = value
		value = p.currentToken.Literal
		p.nextToken()
	}

	expression.Key = key
	expression.Value = value

	if !p.currentTokenIs(token.IN) {
		return nil
	}

	p.nextToken()

	expression.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Block = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		// else if
		if p.peekTokenIs(token.IF) {
			p.nextToken()

			expression.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: p.parseIfExpression(),
					},
				},
			}

			return expression
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseImportExpression() ast.Expression {
	expression := &ast.ImportExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Name = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	p.previousIndexExpression = expression
	p.previousPropertyExpression = nil

	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	while := &ast.WhileExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	while.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	while.Consequence = p.parseBlockStatement()

	return while
}

// ----------------------------------------------------------------------------
// Statements

func (p *Parser) parseAssignStatement() ast.Statement {
	statement := &ast.AssignStatement{}

	if p.currentTokenIs(token.IDENTIFIER) {
		statement.Name = &ast.IdentifierLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	} else if p.currentTokenIs(token.ASSIGN) {
		statement.Token = p.currentToken

		if p.previousIndexExpression != nil {
			statement.Index = p.previousIndexExpression
			p.nextToken()
			statement.Value = p.parseExpression(LOWEST)
			p.previousIndexExpression = nil

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			if fl, ok := statement.Value.(*ast.FunctionLiteral); ok {
				fl.Name = statement.Name.Value
			}

			return statement
		}

		if p.previousPropertyExpression != nil {
			statement.Property = p.previousPropertyExpression
			p.nextToken()
			statement.Value = p.parseExpression(LOWEST)
			p.previousPropertyExpression = nil

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			if fl, ok := statement.Value.(*ast.FunctionLiteral); ok {
				fl.Name = statement.Name.Value
			}

			return statement
		}
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	statement.Token = p.currentToken
	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if fl, ok := statement.Value.(*ast.FunctionLiteral); ok {
		fl.Name = statement.Name.Value
	}

	return statement
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}
