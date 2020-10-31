package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/token"
)

var precedences = map[token.TokenType]int{
	token.BIND:           ASSIGN,
	token.ASSIGN:         ASSIGN,
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
	ASSIGN
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
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

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

	p.registerInfix(token.BIND, p.parseBindExpression)
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
	p.registerInfix(token.DOT, p.parseDotNotationExpression)
	p.registerInfix(token.RANGE, p.parseInfixExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.PLUSPLUS, p.parsePostfixExpression)
	p.registerPostfix(token.MINUSMINUS, p.parsePostfixExpression)

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
	if p.peekToken.Type == token.ASSIGN {
		return p.parseAssignmentStatement()
	}

	switch p.currentToken.Type {
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
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
