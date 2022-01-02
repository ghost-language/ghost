package modules

import (
	"fmt"
	"os"
	"time"

	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

var Os = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Os, "args", osArgs)
	RegisterMethod(Os, "clock", osClock)
	RegisterMethod(Os, "exit", osExit)
	RegisterMethod(Os, "sleep", osSleep)
}

func osArgs(env *object.Environment, args ...object.Object) object.Object {
	list := &object.List{}
	arguments := os.Args[1:]

	for _, argument := range arguments {
		list.Elements = append(list.Elements, &object.String{Value: argument})
	}

	return list
}

func osClock(env *object.Environment, args ...object.Object) object.Object {
	seconds := decimal.NewFromInt(time.Now().Unix())

	return &object.Number{Value: seconds}
}

func osExit(env *object.Environment, args ...object.Object) object.Object {
	var message string

	if len(args) == 2 {
		if args[0].Type() != object.NUMBER {
			// error
			return nil
		}

		if args[1].Type() != object.STRING {
			// error
			return nil
		}

		message = args[1].(*object.String).Value
	} else if len(args) == 1 {
		if args[0].Type() != object.NUMBER {
			// error
			return nil
		}
	} else {
		// error
		return nil
	}

	if message != "" {
		fmt.Println(message)
	}

	arg := args[0].(*object.Number)

	os.Exit(int(arg.Value.IntPart()))

	return arg
}

func osSleep(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		// TODO: error
		return nil
	}

	if args[0].Type() != object.NUMBER {
		// TODO: error
		return nil
	}

	ms := args[0].(*object.Number)
	time.Sleep(time.Duration(ms.Value.IntPart()) * time.Millisecond)

	return nil
}
