package ast

type Node interface {
	Visitable
}

type StatementNode interface {
	Node
}

type ExpressionNode interface {
	Node
}

type Visitor interface {
	visitAssign(*Assign)
	visitBlock(*Block)
	visitBoolean(*Boolean)
	visitCall(*Call)
	visitExpression(*Expression)
	visitFunction(*Function)
	visitIdentifier(*Identifier)
	visitIf(*If)
	visitInfix(*Infix)
	visitList(*List)
	visitNull(*Null)
	visitNumber(*Number)
	visitPrefix(*Prefix)
	visitProgram(*Program)
	visitString(*String)
}

type Visitable interface {
	Accept(Visitor)
}
