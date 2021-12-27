package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) listLiteral() ast.ExpressionNode {
	list := &ast.List{Token: parser.currentToken}

	list.Elements = parser.parseExpressionList(token.RIGHTBRACKET)

	return list
}
