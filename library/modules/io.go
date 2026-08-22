package modules

import (
	"os"
	"path"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var IoMethods = map[string]*object.LibraryFunction{}
var IoProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(IoMethods, "append", ioAppend)
	RegisterMethod(IoMethods, "read", ioRead)
	RegisterMethod(IoMethods, "write", ioWrite)
}

func ioAppend(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("io.append", tok, args, 2); err != nil {
		return err
	}

	target, err := stringAt("io.append", tok, args, 0)

	if err != nil {
		return err
	}

	content, err := stringAt("io.append", tok, args, 1)

	if err != nil {
		return err
	}

	resolved := resolvePath(scope, target)

	file, failure := os.OpenFile(resolved, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if failure != nil {
		return systemFailure("io.append", tok, failure)
	}

	defer file.Close()

	if _, failure := file.WriteString(content + "\n"); failure != nil {
		return systemFailure("io.append", tok, failure)
	}

	return nil
}

func ioRead(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("io.read", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("io.read", tok, args, 0)

	if err != nil {
		return err
	}

	content, failure := os.ReadFile(resolvePath(scope, target))

	if failure != nil {
		return systemFailure("io.read", tok, failure)
	}

	return &object.String{Value: string(content)}
}

func ioWrite(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("io.write", tok, args, 2); err != nil {
		return err
	}

	target, err := stringAt("io.write", tok, args, 0)

	if err != nil {
		return err
	}

	content, err := stringAt("io.write", tok, args, 1)

	if err != nil {
		return err
	}

	resolved := resolvePath(scope, target)

	// The file keeps the permissions it already has, which means it has to
	// already exist. Creating one is `io.append`'s job.
	info, failure := os.Stat(resolved)

	if failure != nil {
		return systemFailure("io.write", tok, failure).
			WithHelp("`io.write` replaces the contents of a file that already exists; use `io.append` to create one")
	}

	if failure := os.WriteFile(resolved, []byte(content), info.Mode()); failure != nil {
		return systemFailure("io.write", tok, failure)
	}

	return nil
}

// resolvePath reads a path as the script would mean it: relative to the
// directory the script itself lives in, not to wherever Ghost was run from.
func resolvePath(scope *object.Scope, target string) string {
	return path.Clean(scope.Environment.GetDirectory() + "/" + target)
}
