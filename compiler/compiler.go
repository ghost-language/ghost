package compiler

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/code"
	"ghostlang.org/x/ghost/object"
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (compiler *Compiler) Compile(node ast.Node) error {
	return nil
}

func (compiler *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}
