package builtins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
	"github.com/shopspring/decimal"
)

func init() {
	RegisterFunction("exit", exitFunction)
	RegisterFunction("first", firstFunction)
	RegisterFunction("identifiers", identifiersFunction)
	RegisterFunction("input", inputFunction)
	RegisterFunction("last", lastFunction)
	RegisterFunction("number", numberFunction)
	RegisterFunction("print", printFunction)
	RegisterFunction("push", pushFunction)
	RegisterFunction("sleep", sleepFunction)
	RegisterFunction("string", stringFunction)
	RegisterFunction("tail", tailFunction)
	RegisterFunction("type", typeFunction)
	RegisterFunction("write", writeFunction)
}

func exitFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	var err object.Object
	var message string

	if len(args) == 2 {
		if args[0].Type() != object.NUMBER_OBJ {
			err = error.NewError(line, error.ArgumentMustBe, "first", "exit", "NUMBER", args[0].Type())
		} else if args[1].Type() != object.STRING_OBJ {
			err = error.NewError(line, error.ArgumentMustBe, "second", "exit", "STRING", args[1].Type())
		}

		message = args[1].(*object.String).Value
	} else if len(args) == 1 {
		if args[0].Type() != object.NUMBER_OBJ {
			err = error.NewError(line, error.ArgumentMustBe, "first", "exit", "NUMBER", args[0].Type())
		}
	} else {
		err = error.NewError(line, error.WrongNumberArguments, len(args), 2)
	}

	if err != nil {
		return err
	}

	if message != "" {
		fmt.Println(message)
	}

	arg := args[0].(*object.Number)
	os.Exit(int(arg.Value.IntPart()))

	return arg
}

func identifiersFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	identifiers := []object.Object{}

	store := env.All()

	for identifier := range store {
		identifiers = append(identifiers, &object.String{Value: identifier})
	}

	return &object.List{Elements: identifiers}
}

func firstFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() != object.LIST_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "first", "LIST", args[0].Type())
	}

	list := args[0].(*object.List)

	return list.Elements[0]
}

func inputFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) == 1 {
		prompt := args[0].(*object.String).Value + " "
		fmt.Fprintf(os.Stdout, prompt)
	}

	buffer := bufio.NewReader(os.Stdin)

	value, _, err := buffer.ReadLine()

	if err != nil && err != io.EOF {
		return error.NewError(line, error.ErrorReadingInput, err)
	}

	return &object.String{Value: string(value)}
}

func lastFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() != object.LIST_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "last", "LIST", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	return list.Elements[length-1]
}

func numberFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() == object.STRING_OBJ {
		arg := args[0].(*object.String)
		num, err := decimal.NewFromString(arg.Value)

		if err != nil {
			return error.NewError(line, error.ArgumentMustBe, "first", "number", "a STRING representation of a number", arg.Value)
		}

		return &object.Number{Value: num}
	}

	if args[0].Type() == object.NUMBER_OBJ {
		return args[0].(*object.Number)
	}

	return error.NewError(line, error.ArgumentMustBe, "first", "number", "a STRING or NUMBER", args[0].Type())
}

func printFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) > 0 {
		fmt.Println(args[0].Inspect())
	} else {
		fmt.Println()
	}

	return nil
}

func writeFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) > 0 {
		fmt.Print(args[0].Inspect())
	} else {
		fmt.Print()
	}

	return nil
}

func pushFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 2 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 2)
	}

	if args[0].Type() != object.LIST_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "push", "LIST", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, list.Elements)
	newElements[length] = args[1]

	list.Elements = newElements

	return &object.Number{Value: decimal.NewFromInt(int64(len(list.Elements)))}
}

func sleepFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() != object.NUMBER_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "sleep", "NUMBER", args[0].Type())
	}

	ms := args[0].(*object.Number)
	time.Sleep(time.Duration(ms.Value.IntPart()) * time.Millisecond)

	return value.NULL
}

func stringFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	return &object.String{Value: args[0].Inspect()}
}

func tailFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() != object.LIST_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "tail", "LIST", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, list.Elements[1:length])

		return &object.List{Elements: newElements}
	}

	return value.NULL
}

func typeFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	val := string(args[0].Type())

	return &object.String{Value: strings.ToLower(val)}
}
