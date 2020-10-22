package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

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
