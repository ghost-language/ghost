package main

import (
	"bytes"
	"syscall/js"

	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
)

func runCode(this js.Value, i []js.Value) interface{} {
	m := make(map[string]interface{})
	var buf bytes.Buffer

	code := i[0].String()
	env := object.NewEnvironment()
	env.SetWriter(&buf)
	l := lexer.New(code)
	p := parser.New(l)

	program := p.ParseProgram()

	result := evaluator.Eval(program, env)

	m["out"] = buf.String()
	m["result"] = result.Inspect()

	return js.ValueOf(m)
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("ghost_run_code", js.FuncOf(runCode))
	<-c
}