package object

import "ghostlang.org/x/ghost/token"

// LibraryModule objects consist of a slice of LibraryFunctions. Classes are
// registered members too, alongside Methods and Properties — `math` has none,
// but a module is free to export a NativeClass the same way it exports a
// method or a property (see library.RegisterClassForScheme).
type LibraryModule struct {
	Name       string
	Methods    map[string]*LibraryFunction
	Properties map[string]*LibraryProperty
	Classes    map[string]*NativeClass
}

// String represents the library module's value as a string.
func (libraryModule *LibraryModule) String() string {
	return libraryModule.Name
}

// Type returns the library module object type.
func (libraryModule *LibraryModule) Type() Type {
	return LIBRARY_MODULE
}

// Method defines the set of methods available on library module objects.
func (libraryModule *LibraryModule) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}
