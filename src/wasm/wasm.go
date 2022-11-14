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
	scope := object.NewScope(object.NewEnvironment())
	scope.Environment.SetWriter(&buf)

	program := parser.Parse(scanner.Scan(code, "wasm.ghost"))

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
