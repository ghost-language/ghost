// Package ast declares the types used to represent syntax trees for
// Ghost source code.
package ast

import (
	"bytes"
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/decimal"
	"ghostlang.org/x/ghost/token"
)

// ----------------------------------------------------------------------------
// Program

// Program is the root node. All programs consist of a slice of
// Statement(s).
type Program struct {
	Statements []Statement
}

// TokenLiteral prints the literal value of the token associated with this node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// String returns a stringified version of the AST for debugging purposes.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ----------------------------------------------------------------------------
// Interfaces

// There are 2 main classes of nodes: Expression nodes, and
// statement nodes.

// Node interface is implemented by all node types.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement interface is implemented by all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression interface is implemented by all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// ----------------------------------------------------------------------------
// Expressions

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.

type (
	// BindExpression defines a new expression type for defining bind expressions.
	BindExpression struct {
		Token token.Token
		Left  Expression
		Value Expression
	}

	// CallExpression defines a new expression type for defining call expressions.
	CallExpression struct {
		Token     token.Token
		Callable  Expression
		Arguments []Expression
	}

	// ForExpression defines a new expression type for defining for expressions.
	// for (x := 0; x <= 10; x += 1) { ... }
	ForExpression struct {
		Token       token.Token
		Identifier  string          // x
		Initializer Statement       // x := 0
		Condition   Expression      // x <= 10
		Increment   Statement       // x += 1
		Block       *BlockStatement // { ... }
	}

	// IfExpression defines a new expression type for defining if expressions.
	IfExpression struct {
		Token       token.Token
		Condition   Expression
		Consequence *BlockStatement
		Alternative *BlockStatement
	}

	// ImportExpression defines a new expression type for defining import expressions.
	ImportExpression struct {
		Token token.Token
		Name  Expression
	}

	// IndexExpression defines a new expression type for defining index expressions.
	IndexExpression struct {
		Token token.Token
		Left  Expression
		Index Expression
	}

	// InfixExpression defines a new expression type for defining infix expressions.
	InfixExpression struct {
		Token    token.Token
		Left     Expression
		Operator string
		Right    Expression
	}

	// PostfixExpression defines a new expression type for defining postfix expressions.
	PostfixExpression struct {
		Token    token.Token
		Operator string
	}

	// PrefixExpression defines a new expression type for defining prefix expressions.
	PrefixExpression struct {
		Token    token.Token
		Operator string
		Right    Expression
	}

	// WhileExpression defines a new expression type for defining while expressions.
	WhileExpression struct {
		Token       token.Token
		Condition   Expression
		Consequence *BlockStatement
	}
)

// ----------------------------------------------------------------------------
// Literals

// A literal is represented by a tree consisting of one
// or more of the following concrete literal nodes.

type (
	BooleanLiteral struct {
		Token token.Token
		Value bool
	}

	IdentifierLiteral struct {
		Token token.Token
		Value string
	}

	ListLiteral struct {
		Token    token.Token
		Elements []Expression
	}

	FunctionLiteral struct {
		Token      token.Token
		Name       string
		Parameters []*IdentifierLiteral
		Defaults   map[string]Expression
		Body       *BlockStatement
	}

	MapLiteral struct {
		Token token.Token
		Pairs map[Expression]Expression
	}

	NumberLiteral struct {
		Token token.Token
		Value decimal.Decimal
	}

	StringLiteral struct {
		Token token.Token
		Value string
	}
)

// expressionNode() ensures that only expression/literal nodes
// can be assigned to an Expression.
//
func (be *BindExpression) expressionNode()    {}
func (ce *CallExpression) expressionNode()    {}
func (fe *ForExpression) expressionNode()     {}
func (ie *IfExpression) expressionNode()      {}
func (ie *ImportExpression) expressionNode()  {}
func (ie *IndexExpression) expressionNode()   {}
func (ie *InfixExpression) expressionNode()   {}
func (pe *PostfixExpression) expressionNode() {}
func (pe *PrefixExpression) expressionNode()  {}
func (we *WhileExpression) expressionNode()   {}

func (bl *BooleanLiteral) expressionNode()    {}
func (fl *FunctionLiteral) expressionNode()   {}
func (il *IdentifierLiteral) expressionNode() {}
func (ll *ListLiteral) expressionNode()       {}
func (ml *MapLiteral) expressionNode()        {}
func (nl *NumberLiteral) expressionNode()     {}
func (sl *StringLiteral) expressionNode()     {}

// TokenLiteral and String implementations for expression/literal nodes.
//
func (be *BindExpression) TokenLiteral() string    { return be.Token.Literal }
func (ce *CallExpression) TokenLiteral() string    { return ce.Token.Literal }
func (fe *ForExpression) TokenLiteral() string     { return fe.Token.Literal }
func (ie *IfExpression) TokenLiteral() string      { return ie.Token.Literal }
func (ie *ImportExpression) TokenLiteral() string  { return ie.Token.Literal }
func (ie *IndexExpression) TokenLiteral() string   { return ie.Token.Literal }
func (ie *InfixExpression) TokenLiteral() string   { return ie.Token.Literal }
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) TokenLiteral() string  { return pe.Token.Literal }
func (we *WhileExpression) TokenLiteral() string   { return we.Token.Literal }

func (bl *BooleanLiteral) TokenLiteral() string    { return bl.Token.Literal }
func (fl *FunctionLiteral) TokenLiteral() string   { return fl.Token.Literal }
func (il *IdentifierLiteral) TokenLiteral() string { return il.Token.Literal }
func (ll *ListLiteral) TokenLiteral() string       { return ll.Token.Literal }
func (ml *MapLiteral) TokenLiteral() string        { return ml.Token.Literal }
func (nl *NumberLiteral) TokenLiteral() string     { return nl.Token.Literal }
func (sl *StringLiteral) TokenLiteral() string     { return sl.Token.Literal }

func (be *BindExpression) String() string {
	var out bytes.Buffer

	out.WriteString(be.Left.String())
	out.WriteString(" " + be.TokenLiteral() + " ")
	out.WriteString(be.Value.String())

	return out.String()
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Callable.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for ")

	out.WriteString(fe.Initializer.String())
	out.WriteString(";")
	out.WriteString(fe.Condition.String())
	out.WriteString(";")
	out.WriteString(fe.Increment.String())
	out.WriteString(";")
	out.WriteString(fe.Block.String())

	return out.String()
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (ie *ImportExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral())
	out.WriteString("(")
	out.WriteString(fmt.Sprintf("\"%s\"", ie.Name))
	out.WriteString(")")

	return out.String()
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Token.Literal)
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Consequence.String())

	return out.String()
}

func (bl *BooleanLiteral) String() string { return bl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, p := range fl.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

func (il *IdentifierLiteral) String() string { return il.Value }

func (ll *ListLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (ml *MapLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for key, value := range ml.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (nl *NumberLiteral) String() string { return nl.Token.Literal }
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.

type (
	// AssignmentStatement defines a new statement type for defining assignments.
	AssignmentStatement struct {
		Token token.Token
		Name  *IdentifierLiteral
		Value Expression
	}

	// BlockStatement defines a new statement type for defining blocks.
	BlockStatement struct {
		Token      token.Token
		Statements []Statement
	}

	// ExpressionStatement defines a new statement type for defining expressions.
	ExpressionStatement struct {
		Token      token.Token
		Expression Expression
	}

	// ReturnStatement defines a new statement type for defining returns.
	ReturnStatement struct {
		Token       token.Token
		ReturnValue Expression
	}
)

// statementNode() ensures that only statement nodes
// can be assigned to a Statement.
//
func (as *AssignmentStatement) statementNode() {}
func (bs *BlockStatement) statementNode()      {}
func (es *ExpressionStatement) statementNode() {}
func (rs *ReturnStatement) statementNode()     {}

// TokenLiteral and String implementations for statement nodes.
//
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (bs *BlockStatement) TokenLiteral() string      { return bs.Token.Literal }
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (rs *ReturnStatement) TokenLiteral() string     { return rs.Token.Literal }

func (as *AssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(as.TokenLiteral() + " ")
	out.WriteString(as.Value.String())

	return out.String()
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
