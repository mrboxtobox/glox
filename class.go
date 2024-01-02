package main

import "fmt"

// Runtime representation of a glox class.
type Class struct {
	// Name of the class.
	Name string
	// Optional superclass extended by the class.
	Superclass *Class
	// Methods available to the class.
	Methods map[string]Function
}

// Creates a new class.
func NewClass(name string, superclass *Class, methods map[string]Function) *Class {
	return &Class{name, superclass, methods}
}

// Searches for a method in a class or its inheritance chain.
// Returns the method and a boolean indicating if the method was found.
func (c *Class) FindMethod(name string) (Function, bool) {
	if method, found := c.Methods[name]; found {
		return method, true
	}
	if c.Superclass != nil {
		return c.Superclass.FindMethod(name)
	}
	return Function{}, false
}

// Number of arguments used in the initializer, if present. Otherwise, it is 0.
func (c *Class) Arity() int {
	if initializer, found := c.FindMethod("init"); found {
		return initializer.Arity()
	}
	return 0
}

// Creates a new instance and runs the initializer, if an initializer exists.
// Returns the new instance and an error (if any) from initialization.
func (c *Class) Call(interpreter Interpreter, arguments []any) (any, error) {
	instance := NewInstance(c)
	if initializer, found := c.FindMethod("init"); found {
		if _, err := initializer.Bind(instance).Call(interpreter, arguments); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (c *Class) String() string {
	return fmt.Sprintf("<class %s>", c.Name)
}
