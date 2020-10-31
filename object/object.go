// Package object describes a named language entity such as
// a type, variable, function, or literal. All objects
// implement the Object interface.
package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/decimal"
)

// ----------------------------------------------------------------------------
// Constants

const (
	BOOLEAN_OBJ      = "BOOLEAN"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	LIST_OBJ         = "LIST"
	MAP_OBJ          = "MAP"
	MODULE_OBJ       = "MODULE"
	NULL_OBJ         = "NULL"
	NUMBER_OBJ       = "NUMBER"
	PACKAGE_OBJ      = "PACKAGE"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	STRING_OBJ       = "STRING"
)

// ----------------------------------------------------------------------------
// Interfaces

type ObjectType string

// Object interface is implemented by all objects.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Mutable interface is implemented by all mutable objects.
type Mutable interface {
	Set(obj Object)
}

type (
	Boolean struct {
		Value bool
	}

	Builtin struct {
		Name string
		Fn   BuiltinFunction
	}

	BuiltinFunction func(env *Environment, args ...Object) Object

	Module struct {
		Name      string
		Functions []BuiltinFunction
	}

	Error struct {
		Message string
	}

	Function struct {
		Parameters []*ast.IdentifierLiteral
		Body       *ast.BlockStatement
		Defaults   map[string]ast.Expression
		Env        *Environment
	}

	List struct {
		Elements []Object
	}

	Map struct {
		Pairs map[MapKey]MapPair
	}

	Mappable interface {
		MapKey() MapKey
	}

	// MapKey defines the key for maps that can be comparable and unique.
	MapKey struct {
		Type  ObjectType
		Value uint64
	}

	MapPair struct {
		Key   Object
		Value Object
	}

	Null struct{}

	Number struct {
		Value decimal.Decimal
	}

	Package struct {
		Name       string
		Attributes Object
	}

	ReturnValue struct {
		Value Object
	}

	String struct {
		Value string
	}
)

// ----------------------------------------------------------------------------
// Types

func (b *Boolean) Type() ObjectType      { return BOOLEAN_OBJ }
func (b *Builtin) Type() ObjectType      { return BUILTIN_OBJ }
func (e *Error) Type() ObjectType        { return ERROR_OBJ }
func (f *Function) Type() ObjectType     { return FUNCTION_OBJ }
func (l *List) Type() ObjectType         { return LIST_OBJ }
func (m *Map) Type() ObjectType          { return MAP_OBJ }
func (m *Module) Type() ObjectType       { return MODULE_OBJ }
func (n *Null) Type() ObjectType         { return NULL_OBJ }
func (n *Number) Type() ObjectType       { return NUMBER_OBJ }
func (m *Package) Type() ObjectType      { return PACKAGE_OBJ }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (s *String) Type() ObjectType       { return STRING_OBJ }

// ----------------------------------------------------------------------------
// Inspections

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Builtin) Inspect() string { return fmt.Sprintf("builtin function: %s", b.Name) }
func (e *Error) Inspect() string   { return "ERROR: " + e.Message }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString("function")
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String() + "\n")
	out.WriteString("}\n")

	return out.String()
}

func (l *List) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range l.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (m *Map) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (m *Module) Inspect() string       { return fmt.Sprintf("module(%s)", m.Name) }
func (n *Null) Inspect() string         { return "null" }
func (n *Number) Inspect() string       { return n.Value.String() }
func (m *Package) Inspect() string      { return fmt.Sprintf("package(%s)", m.Name) }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }
func (s *String) Inspect() string       { return s.Value }

// ----------------------------------------------------------------------------
// Settables

func (b *Boolean) Set(obj Object) { b.Value = obj.(*Boolean).Value }
func (e *Error) Set(obj Object)   { e.Message = obj.(*Error).Message }
func (l *List) Set(obj Object)    { l.Elements = obj.(*List).Elements }
func (n *Number) Set(obj Object)  { n.Value = obj.(*Number).Value }
func (s *String) Set(obj Object)  { s.Value = obj.(*String).Value }

// ----------------------------------------------------------------------------
// Mappables

func (b *Boolean) MapKey() MapKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return MapKey{Type: b.Type(), Value: value}
}

func (n *Number) MapKey() MapKey {
	value, _ := strconv.ParseUint(n.Value.String(), 10, 64)

	return MapKey{Type: n.Type(), Value: value}
}

func (s *String) MapKey() MapKey {
	// Note: There is a _very_ small chance that the following will
	// result in the same hash being generated for different string
	// values (hash collisions). Research "separate chaining" and
	// "open addressing" techniques to work around the problem.

	hash := fnv.New64a()
	hash.Write([]byte(s.Value))

	return MapKey{Type: s.Type(), Value: hash.Sum64()}
}
