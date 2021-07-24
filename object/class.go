package object

import "fmt"

// CLASS represents the object's type.
const CLASS = "CLASS"

type Class struct {
	// Callable
	Name string
	Methods map[string]*UserFunction
}

func (c *Class) String() string {
	return fmt.Sprintf("<class %s>", c.Name)
}

// Type returns the class object type.
func (c *Class) Type() Type {
	return CLASS
}

// func (c *Class) Call(arguments []Object) (Object, error) {
// 	instance := &ClassInstance{Class: c, fields: make(map[string]interface{})}

// 	if constructor, ok := c.Methods["constructor"]; ok {
// 		constructor.Bind(instance).Call(arguments)

// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 	}

// 	return instance, nil
// }

func (c *Class) Artity() int {
	if constructor, ok := c.Methods["constructor"]; ok {
		return constructor.Arity()
	}

	return 0
}