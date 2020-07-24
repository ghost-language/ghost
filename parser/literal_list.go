package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (p *Parser) parseListLiteral() ast.Expression {
	list := &ast.ListLiteral{Token: p.currentToken}
	list.Elements = p.parseExpressionList(token.RBRACKET)

	return list
}
