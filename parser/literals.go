package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/decimal"
	"ghostlang.org/x/ghost/token"
)

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
