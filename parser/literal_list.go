package parser

import (
	"ghostlang.org/ghost/ast"
	"ghostlang.org/ghost/token"
)

func (p *Parser) parseListLiteral() ast.Expression {
	list := &ast.ListLiteral{Token: p.currentToken}
	list.Elements = p.parseExpressionList(token.RBRACKET)

	return list
}
