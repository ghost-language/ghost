package parser

import "ghostlang.org/x/ghost/ast"

// spreadExpression parses `...expr` - see ast.Spread. Its operand binds at
// prefix precedence, the same as `!`/`-`, so it reaches only as far as the
// next comma in the surrounding call or list literal.
func (parser *Parser) spreadExpression() ast.ExpressionNode {
	spread := &ast.Spread{Token: parser.currentToken}

	parser.readToken()

	spread.Value = parser.parseExpression(PREFIX)

	return spread
}
