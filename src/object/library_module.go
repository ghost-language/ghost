package object

import "fmt"

const LIBRARY_MODULE = "LIBRARY_MODULE"

// LibraryModule objects consist of a slice of LibraryFunctions.
type LibraryModule struct {
	Name    string
	Methods map[string]*LibraryFunction
}

// String represents the library module's value as a string.
func (libraryModule *LibraryModule) String() string {
	return fmt.Sprintf("library module {%s}", libraryModule.Name)
}

// Type returns the library module object type.
func (libraryModule *LibraryModule) Type() Type {
	return LIBRARY_MODULE
}

// Method defines the set of methods available on library module objects.
func (libraryModule *LibraryModule) Method(method string, args []Object) (Object, bool) {
	return nil, false
}