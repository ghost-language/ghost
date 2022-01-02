package modules

import (
	"io/ioutil"
	"os"
	"path"

	"ghostlang.org/x/ghost/object"
)

var Io = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Io, "append", ioAppend)
	RegisterMethod(Io, "read", ioRead)
	RegisterMethod(Io, "write", ioWrite)
}

func ioAppend(env *object.Environment, args ...object.Object) object.Object {
	path := path.Clean(env.GetDirectory() + "/" + args[0].(*object.String).Value)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil
	}

	defer file.Close()

	file.WriteString(args[1].(*object.String).Value + "\n")

	return nil
}

func ioRead(env *object.Environment, args ...object.Object) object.Object {
	path := path.Clean(env.GetDirectory() + "/" + args[0].(*object.String).Value)
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil
	}

	return &object.String{Value: string(content)}
}

func ioWrite(env *object.Environment, args ...object.Object) object.Object {
	path := path.Clean(env.GetDirectory() + "/" + args[0].(*object.String).Value)
	contents := []byte(args[1].(*object.String).Value)
	info, _ := os.Stat(path)
	mode := info.Mode()

	err := ioutil.WriteFile(path, contents, mode)

	if err != nil {
		return nil
	}

	return nil
}
