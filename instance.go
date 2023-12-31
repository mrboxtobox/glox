package main

type Instance struct {
	Class *Class
	// Struct-internal.
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		Class:  class,
		fields: map[string]any{},
	}
}

func (i *Instance) Get(name Token) (any, error) {
	if object, found := i.fields[name.Lexeme]; found {
		return object, nil
	}
	if method, found := i.Class.FindMethod(name.Lexeme); found {
		return method.Bind(i), nil
	}
	err := LogAndReturnError(name, "Undefined property '"+name.Lexeme+"'.")
	return nil, err
}

func (i *Instance) Set(name Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *Instance) String() string {
	return i.Class.Name + " instance"
}
