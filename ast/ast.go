package ast

type Node interface {
	Accept(Visitor)
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
	visitIdentifier(*Identifier)
	visitIf(*If)
	visitInfix(*Infix)
	visitNull(*Null)
	visitNumber(*Number)
	visitPrefix(*Prefix)
	visitProgram(*Program)
	visitString(*String)
}
