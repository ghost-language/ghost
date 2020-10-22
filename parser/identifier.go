package parser

import "ghostlang.org/x/ghost/ast"

func (p *Parser) parseIdentifierLiteral() ast.Expression {
	return &ast.IdentifierLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}
