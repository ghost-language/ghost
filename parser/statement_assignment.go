package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	statement := &ast.AssignmentStatement{Token: p.peekToken}
	statement.Name = &ast.IdentifierLiteral{Token: p.currentToken, Value: p.currentToken.Literal}

	p.nextToken()
	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}
