package parser

import "ghostlang.org/x/ghost/ast"

// postfixExpression parses `++` and `--`. It is registered as an infix parser
// so it receives the expression already parsed to its left — an identifier, a
// property, or an index — rather than only the bare name a plain identifier
// would give it. That is what lets `this.score++` and `list[0]++` target the
// property or element instead of a same-named local variable.
func (parser *Parser) postfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	return &ast.Postfix{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Type,
		Left:     left,
	}
}
