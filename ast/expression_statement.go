package ast

import "ghostlang.org/ghost/token"

/*
ExpressionStatement defines a new statement for defining expressions.

It has two fields:

- The Token field, which every node has
- The Expression field, which holds the expression

These two fields fulfills the Statement interface, which means
we can add it to the Statements slice of our AST program.
*/
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the token literal value.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
