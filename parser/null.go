package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) null() ast.ExpressionNode {
	null := &ast.Null{Token: parser.current()}

	return null
}
