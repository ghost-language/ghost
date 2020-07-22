package parser

import (
	"fmt"

	"ghostlang.org/ghost/ast"
)

func (p *Parser) parseBindExpression(expression ast.Expression) ast.Expression {
	switch node := expression.(type) {
	case *ast.Identifier:
	default:
		message := fmt.Sprintf("expected identifier expression on left but got %T (%+v)", node, expression)
		p.errors = append(p.errors, message)
		return nil
	}

	be := &ast.BindExpression{Token: p.currentToken, Left: expression}

	p.nextToken()

	be.Value = p.parseExpression(LOWEST)

	if fl, ok := be.Value.(*ast.FunctionLiteral); ok {
		identifier := be.Left.(*ast.Identifier)
		fl.Name = identifier.Value
	}

	return be
}
