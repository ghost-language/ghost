package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) identifier() ast.ExpressionNode {
	identifier := &ast.Identifier{Token: parser.current(), Value: parser.current().Literal}

	return identifier
}
