package main

import (
	"bytes"
	"syscall/js"

	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func ghost(this js.Value, i []js.Value) interface{} {
	m := make(map[string]interface{})
	var buf bytes.Buffer

	code := i[0].String()
	scope := &object.Scope{Environment: object.NewEnvironment()}
	scope.Environment.SetWriter(&buf)

	scanner := scanner.New(code, "wasm.ghost")
	parser := parser.New(scanner)
	program := parser.Parse()

	result := evaluator.Evaluate(program, scope)

	m["result"] = buf.String()
	m["object"] = result.String()

	return js.ValueOf(m)
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("ghost", js.FuncOf(ghost))
	<-c
}
