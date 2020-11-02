package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (p *Parser) parseAssignStatement() ast.Statement {
	statement := &ast.AssignStatement{}

	if p.currentTokenIs(token.IDENTIFIER) {
		statement.Name = &ast.IdentifierLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	} else if p.currentTokenIs(token.ASSIGN) {
		statement.Token = p.currentToken

		if p.previousIndexExpression != nil {
			statement.Index = p.previousIndexExpression
			p.nextToken()
			statement.Value = p.parseExpression(LOWEST)
			p.previousIndexExpression = nil

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			if fl, ok := statement.Value.(*ast.FunctionLiteral); ok {
				fl.Name = statement.Name.Value
			}

			return statement
		}

		if p.previousPropertyExpression != nil {
			statement.Property = p.previousPropertyExpression
			p.nextToken()
			statement.Value = p.parseExpression(LOWEST)
			p.previousPropertyExpression = nil

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			if fl, ok := statement.Value.(*ast.FunctionLiteral); ok {
				fl.Name = statement.Name.Value
			}

			return statement
		}
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	statement.Token = p.currentToken
	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if fl, ok := statement.Value.(*ast.FunctionLiteral); ok {
		fl.Name = statement.Name.Value
	}

	return statement
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}
