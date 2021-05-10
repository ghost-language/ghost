package builtins

import (
	"fmt"
	"io/ioutil"
	"os"

	"ghostlang.org/x/ghost/object"
)

func init() {
	RegisterFunction("file.read", fileReadFunction)
	RegisterFunction("file.write", fileWriteFunction)
}

// fileReadFunction returns the arguments passed to the program
func fileReadFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	path := env.GetDirectory() + args[0].Inspect()
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return &object.Error{Message: err.Error() + fmt.Sprintf(" on line %d", line)}
	}

	return &object.String{Value: string(content)}
}

func fileWriteFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	path := env.GetDirectory() + args[0].Inspect()
	contents := []byte(args[1].Inspect())
	info, _ := os.Stat(path)
	mode := info.Mode()

	err := ioutil.WriteFile(path, contents, mode)

	if err != nil {
		return &object.Error{Message: err.Error() + fmt.Sprintf(" on line %d", line)}
	}

	return &object.Null{}
}