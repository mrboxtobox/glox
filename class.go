package main

type Class struct {
	Name    string
	Methods map[string]Function
}

func NewClass(name string, methods map[string]Function) *Class {
	return &Class{name, methods}
}

func (c *Class) FindMethod(name string) (Function, bool) {
	if method, found := c.Methods[name]; found {
		return method, true
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
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance, nil
}

func (c *Class) String() string {
	return c.Name
}
