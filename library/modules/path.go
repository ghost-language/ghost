package modules

import (
	"path/filepath"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var PathMethods = map[string]*object.LibraryFunction{}
var PathProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(PathMethods, "basename", pathBasename)
	RegisterMethod(PathMethods, "dirname", pathDirname)
	RegisterMethod(PathMethods, "extname", pathExtname)
	RegisterMethod(PathMethods, "isAbsolute", pathIsAbsolute)
	RegisterMethod(PathMethods, "join", pathJoin)
}

// pathJoin and the rest of this module are pure string manipulation — unlike
// `file`, nothing here touches the filesystem or needs to know where the
// running script lives, which is the actual difference between the two
// modules (§7: a module names a domain, and "building a path" and "reading
// one" are different domains even though scripts usually do both together).
func pathJoin(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityAtLeast("path.join", tok, args, 1); err != nil {
		return err
	}

	parts := make([]string, len(args))

	for index := range args {
		part, err := stringAt("path.join", tok, args, index)

		if err != nil {
			return err
		}

		parts[index] = part
	}

	return &object.String{Value: filepath.Join(parts...)}
}

func pathBasename(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("path.basename", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("path.basename", tok, args, 0)

	if err != nil {
		return err
	}

	return &object.String{Value: filepath.Base(target)}
}

func pathDirname(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("path.dirname", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("path.dirname", tok, args, 0)

	if err != nil {
		return err
	}

	return &object.String{Value: filepath.Dir(target)}
}

func pathExtname(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("path.extname", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("path.extname", tok, args, 0)

	if err != nil {
		return err
	}

	return &object.String{Value: filepath.Ext(target)}
}

func pathIsAbsolute(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("path.isAbsolute", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("path.isAbsolute", tok, args, 0)

	if err != nil {
		return err
	}

	return &object.Boolean{Value: filepath.IsAbs(target)}
}
