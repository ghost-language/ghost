package builtins

import (
	"fmt"
	"net/http"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/server/router"
)

func init() {
	RegisterFunction("Router.get", routerGetFunction)
	RegisterFunction("Router.post", routerPostFunction)
}

// routerGetFunction ...
func routerGetFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	router.RegisterRoute("GET", args[0].Inspect(), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET: %s", args[0].Inspect())
		// function := args[1].(*object.Function)
		// extendedEnv := evaluator.ExtendFunctionEnv(function, []object.Object{})

		// evaluated := evaluator.Eval(function.Body, extendedEnv)

		// evaluator.UnwrapReturnValue(evaluated)
	})

	return &object.Null{}
}

// routerPostFunction ...
func routerPostFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	router.RegisterRoute("POST", args[0].Inspect(), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test route called")
	})

	return &object.Null{}
}