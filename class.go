package main

type Class struct {
	Name       string
	Superclass *Class
	Methods    map[string]Function
}

func NewClass(name string, superclass *Class, methods map[string]Function) *Class {
	return &Class{name, superclass, methods}
}

func (c *Class) FindMethod(name string) (Function, bool) {
	if method, found := c.Methods[name]; found {
		return method, true
	}
	if c.Superclass != nil {
		return c.Superclass.FindMethod(name)
	}
	return Function{}, false
}

func (c *Class) Arity() int {
	if initializer, found := c.FindMethod("init"); found {
		return initializer.Arity()
	}
	return 0
}

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
	return c.Name
}
