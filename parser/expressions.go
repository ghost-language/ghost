package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (p *Parser) parseCallExpression(callable ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Callable: callable}
	expression.Arguments = p.parseExpressionList(token.RPAREN)

	return expression
}

func (p *Parser) parseDotExpression(object ast.Expression) ast.Expression {
	tok := p.currentToken

	precedence := p.currentPrecendence()
	p.nextToken()

	if p.peekTokenIs(token.LPAREN) {
		// method
		exp := &ast.MethodExpression{Token: tok, Object: object}
		exp.Method = p.parseExpression(precedence)
		p.nextToken()
		exp.Arguments = p.parseExpressionList(token.RPAREN)

		return exp
	}

	// property
	exp := &ast.PropertyExpression{Token: tok, Object: object}
	exp.Property = p.parseIdentifierLiteral()

	p.previousPropertyExpression = exp
	p.previousIndexExpression = nil

	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	if !p.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return p.parseForInExpression(expression)
	}

	expression.Identifier = p.currentToken.Literal
	expression.Initializer = p.parseAssignStatement()

	if expression.Initializer == nil {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if expression.Condition == nil {
		return nil
	}

	p.nextToken()
	p.nextToken()

	expression.Increment = p.parseAssignStatement()

	if expression.Increment == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Block = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForInExpression(parentExpression *ast.ForExpression) ast.Expression {
	expression := &ast.ForInExpression{Token: parentExpression.Token}

	if !p.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	value := p.currentToken.Literal
	var key string
	p.nextToken()

	if p.currentTokenIs(token.COMMA) {
		p.nextToken()

		if !p.currentTokenIs(token.IDENTIFIER) {
			return nil
		}

		key = value
		value = p.currentToken.Literal
		p.nextToken()
	}

	expression.Key = key
	expression.Value = value

	if !p.currentTokenIs(token.IN) {
		return nil
	}

	p.nextToken()

	expression.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Block = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		// else if
		if p.peekTokenIs(token.IF) {
			p.nextToken()

			expression.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: p.parseIfExpression(),
					},
				},
			}

			return expression
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseImportExpression() ast.Expression {
	expression := &ast.ImportExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Name = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	p.previousIndexExpression = expression
	p.previousPropertyExpression = nil

	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	while := &ast.WhileExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	while.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	while.Consequence = p.parseBlockStatement()

	return while
}
