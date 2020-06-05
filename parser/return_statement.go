package parser

import (
	"ghostlang.org/ghost/ast"
	"ghostlang.org/ghost/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}
