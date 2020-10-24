package builtins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"ghostlang.org/x/ghost/decimal"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/utilities"
)

func init() {
	RegisterFunction("first", firstFunction)
	RegisterFunction("input", inputFunction)
	RegisterFunction("last", lastFunction)
	RegisterFunction("len", lenFunction)
	RegisterFunction("print", printFunction)
	RegisterFunction("push", pushFunction)
	RegisterFunction("tail", tailFunction)
	RegisterFunction("identifiers", identifiersFunction)
}

func identifiersFunction(env *object.Environment, args ...object.Object) object.Object {
	identifiers := []object.Object{}

	store := env.All()

	for identifier := range store {
		identifiers = append(identifiers, &object.String{Value: identifier})
	}

	return &object.List{Elements: identifiers}
}

func firstFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return utilities.NewError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return utilities.NewError("argument to `first` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)

	return list.Elements[0]
}

func inputFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) == 1 {
		prompt := args[0].(*object.String).Value + " "
		fmt.Fprintf(os.Stdout, prompt)
	}

	buffer := bufio.NewReader(os.Stdin)

	line, _, err := buffer.ReadLine()

	if err != nil && err != io.EOF {
		return utilities.NewError(fmt.Sprintf("error reading input: %s", err))
	}
	return &object.String{Value: string(line)}
}

func lastFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return utilities.NewError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return utilities.NewError("argument to `last` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	return list.Elements[length-1]
}

func lenFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return utilities.NewError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.List:
		return &object.Number{Value: decimal.NewFromInt(int64(len(arg.Elements)))}
	case *object.String:
		return &object.Number{Value: decimal.NewFromInt(int64(utf8.RuneCountInString(arg.Value)))}
	default:
		return utilities.NewError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func printFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) > 0 {
		fmt.Println(args[0].Inspect())
	} else {
		fmt.Println()
	}

	return nil
}

func pushFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return utilities.NewError("wrong number of arguments. got=%d, expected=2", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return utilities.NewError("argument to `push` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, list.Elements)
	newElements[length] = args[1]

	list.Elements = newElements

	return &object.Number{Value: decimal.NewFromInt(int64(len(list.Elements)))}
}

func tailFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return utilities.NewError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return utilities.NewError("argument to `tail` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, list.Elements[1:length])

		return &object.List{Elements: newElements}
	}

	return NULL
}
