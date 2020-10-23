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
	"unicode"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/decimal"
)

// ----------------------------------------------------------------------------
// Environment

// Environment is an object that holds a mapping of names to bound objects
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

// All returns all stored identifiers.
func (e *Environment) All() map[string]Object {
	return e.store
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}

func (e *Environment) Exported() *Map {
	pairs := make(map[MapKey]MapPair)

	for k, v := range e.store {
		// Replace this with checking for "Import" token
		if unicode.IsUpper(rune(k[0])) {
			s := &String{Value: k}
			pairs[s.MapKey()] = MapPair{Key: s, Value: v}
		}
	}

	return &Map{Pairs: pairs}
}

// ----------------------------------------------------------------------------
// Constants

const (
	BOOLEAN_OBJ      = "BOOLEAN"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	LIST_OBJ         = "LIST"
	MAP_OBJ          = "MAP"
	NULL_OBJ         = "NULL"
	NUMBER_OBJ       = "NUMBER"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	STRING_OBJ       = "STRING"
	MODULE_OBJ       = "MODULE"
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

	Module struct {
		Name       string
		Attributes Object
	}

	Null struct{}

	Number struct {
		Value decimal.Decimal
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
func (ao *List) Type() ObjectType        { return LIST_OBJ }
func (m *Map) Type() ObjectType          { return MAP_OBJ }
func (m *Module) Type() ObjectType       { return MODULE_OBJ }
func (n *Null) Type() ObjectType         { return NULL_OBJ }
func (n *Number) Type() ObjectType       { return NUMBER_OBJ }
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
