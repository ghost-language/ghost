package modules

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var FileMethods = map[string]*object.LibraryFunction{}
var FileProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(FileMethods, "append", fileAppend)
	RegisterMethod(FileMethods, "copy", fileCopy)
	RegisterMethod(FileMethods, "delete", fileDelete)
	RegisterMethod(FileMethods, "exists", fileExists)
	RegisterMethod(FileMethods, "isDirectory", fileIsDirectory)
	RegisterMethod(FileMethods, "list", fileList)
	RegisterMethod(FileMethods, "mkdir", fileMkdir)
	RegisterMethod(FileMethods, "move", fileMove)
	RegisterMethod(FileMethods, "read", fileRead)
	RegisterMethod(FileMethods, "size", fileSize)
	RegisterMethod(FileMethods, "write", fileWrite)
}

func fileAppend(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.append", tok, args, 2); err != nil {
		return err
	}

	target, err := stringAt("file.append", tok, args, 0)

	if err != nil {
		return err
	}

	content, err := stringAt("file.append", tok, args, 1)

	if err != nil {
		return err
	}

	resolved := resolvePath(scope, target)

	file, failure := os.OpenFile(resolved, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if failure != nil {
		return systemFailure("file.append", tok, failure)
	}

	defer file.Close()

	if _, failure := file.WriteString(content + "\n"); failure != nil {
		return systemFailure("file.append", tok, failure)
	}

	return nil
}

func fileRead(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.read", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.read", tok, args, 0)

	if err != nil {
		return err
	}

	content, failure := os.ReadFile(resolvePath(scope, target))

	if failure != nil {
		return systemFailure("file.read", tok, failure)
	}

	return &object.String{Value: string(content)}
}

func fileWrite(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.write", tok, args, 2); err != nil {
		return err
	}

	target, err := stringAt("file.write", tok, args, 0)

	if err != nil {
		return err
	}

	content, err := stringAt("file.write", tok, args, 1)

	if err != nil {
		return err
	}

	resolved := resolvePath(scope, target)

	// The file keeps the permissions it already has, which means it has to
	// already exist. Creating one is `file.append`'s job.
	info, failure := os.Stat(resolved)

	if failure != nil {
		return systemFailure("file.write", tok, failure).
			WithHelp("`file.write` replaces the contents of a file that already exists; use `file.append` to create one")
	}

	if failure := os.WriteFile(resolved, []byte(content), info.Mode()); failure != nil {
		return systemFailure("file.write", tok, failure)
	}

	return nil
}

func fileExists(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.exists", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.exists", tok, args, 0)

	if err != nil {
		return err
	}

	_, failure := os.Stat(resolvePath(scope, target))

	return &object.Boolean{Value: failure == nil}
}

func fileIsDirectory(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.isDirectory", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.isDirectory", tok, args, 0)

	if err != nil {
		return err
	}

	info, failure := os.Stat(resolvePath(scope, target))

	if failure != nil {
		return systemFailure("file.isDirectory", tok, failure)
	}

	return &object.Boolean{Value: info.IsDir()}
}

func fileSize(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.size", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.size", tok, args, 0)

	if err != nil {
		return err
	}

	info, failure := os.Stat(resolvePath(scope, target))

	if failure != nil {
		return systemFailure("file.size", tok, failure)
	}

	return object.NewInt(info.Size())
}

func fileDelete(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.delete", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.delete", tok, args, 0)

	if err != nil {
		return err
	}

	// Removes a file or an empty directory — the same reach as Go's own
	// os.Remove. A non-empty directory is left alone rather than silently
	// wiped, which is `file.delete`'s job to refuse, not the caller's to
	// guard against.
	if failure := os.Remove(resolvePath(scope, target)); failure != nil {
		return systemFailure("file.delete", tok, failure)
	}

	return nil
}

func fileMkdir(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.mkdir", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.mkdir", tok, args, 0)

	if err != nil {
		return err
	}

	// Creates any missing parent directories too, like `mkdir -p` — a script
	// asking for a directory to exist almost never wants to be told one of
	// its ancestors didn't.
	if failure := os.MkdirAll(resolvePath(scope, target), 0755); failure != nil {
		return systemFailure("file.mkdir", tok, failure)
	}

	return nil
}

func fileList(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.list", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("file.list", tok, args, 0)

	if err != nil {
		return err
	}

	entries, failure := os.ReadDir(resolvePath(scope, target))

	if failure != nil {
		return systemFailure("file.list", tok, failure)
	}

	names := make([]object.Object, len(entries))

	for index, entry := range entries {
		names[index] = &object.String{Value: entry.Name()}
	}

	return &object.List{Elements: names}
}

func fileCopy(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.copy", tok, args, 2); err != nil {
		return err
	}

	source, err := stringAt("file.copy", tok, args, 0)

	if err != nil {
		return err
	}

	destination, err := stringAt("file.copy", tok, args, 1)

	if err != nil {
		return err
	}

	resolvedSource := resolvePath(scope, source)

	info, failure := os.Stat(resolvedSource)

	if failure != nil {
		return systemFailure("file.copy", tok, failure)
	}

	input, failure := os.Open(resolvedSource)

	if failure != nil {
		return systemFailure("file.copy", tok, failure)
	}

	defer input.Close()

	output, failure := os.OpenFile(resolvePath(scope, destination), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

	if failure != nil {
		return systemFailure("file.copy", tok, failure)
	}

	defer output.Close()

	if _, failure := io.Copy(output, input); failure != nil {
		return systemFailure("file.copy", tok, failure)
	}

	return nil
}

func fileMove(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("file.move", tok, args, 2); err != nil {
		return err
	}

	source, err := stringAt("file.move", tok, args, 0)

	if err != nil {
		return err
	}

	destination, err := stringAt("file.move", tok, args, 1)

	if err != nil {
		return err
	}

	if failure := os.Rename(resolvePath(scope, source), resolvePath(scope, destination)); failure != nil {
		return systemFailure("file.move", tok, failure)
	}

	return nil
}

// resolvePath reads a path as the script would mean it: relative to the
// directory the script itself lives in, not to wherever Ghost was run from.
func resolvePath(scope *object.Scope, target string) string {
	if filepath.IsAbs(target) {
		return path.Clean(target)
	}

	return path.Clean(scope.Environment.GetDirectory() + "/" + target)
}
