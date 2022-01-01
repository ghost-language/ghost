package object

import "fmt"

const LIBRARY_FUNCTION = "LIBRARY_FUNCTION"

// LibraryFunction objects consist of a native Go function.
type LibraryFunction struct {
	Name     string
	Function GoFunction
}

func (object *LibraryFunction) Accept(v Visitor) {
	v.visitLibraryFunction(object)
}

// String represents the library function's value as a string.
func (libraryFunction *LibraryFunction) String() string {
	return fmt.Sprintf("library function {%s}", libraryFunction.Name)
}

// Type returns the library function object type.
func (libraryFunction *LibraryFunction) Type() Type {
	return LIBRARY_FUNCTION
}

func (libraryFunction *LibraryFunction) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
