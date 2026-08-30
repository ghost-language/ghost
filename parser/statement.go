package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) statement() ast.StatementNode {
	switch parser.currentToken.Type {
	case token.RETURN:
		return parser.returnStatement()
	}

	statement := parser.assign()

	if statement != nil {
		return statement
	}

	return parser.expressionStatement()
}

func (parser *Parser) expressionStatement() ast.StatementNode {
	return parser.expressionStatementFrom(parser.parseExpression(LOWEST))
}

// expressionStatementFrom wraps an already-parsed expression as a statement,
// consuming an optional trailing `;` the same way assign()/returnStatement()/
// destructuringAssign()'s own successful-pattern branch already do.
// expressionStatement() is the path every bare call, and every
// if/while/for/function/class/trait/switch/import/use/break/continue
// statement, funnels through, as an expression whose prefix parser happens
// to build one of those nodes; destructuringAssign() reaches here too, on
// each of its "this wasn't actually a pattern" fallbacks, since a plain
// list/map literal statement (`[1, 2, 3]`, `{"a": 1}`) starts with the exact
// same tokens as a destructuring pattern. Before this existed, `;` was
// reliable only after a plain assignment - it broke parsing of the very next
// statement after anything else (§13.12).
func (parser *Parser) expressionStatementFrom(expression ast.ExpressionNode) ast.StatementNode {
	statement := &ast.Expression{Expression: expression}

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.readToken()
	}

	return statement
}
