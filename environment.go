package main

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: map[string]any{},
	}
}

func (e *Environment) ToString() string {
	return fmt.Sprintf("Environment{\n\tenclosing=%+v,\nvalues=%v\n}", e.enclosing, e.values)
}

func NewEnvironmentFromEnclosing(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    map[string]any{},
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
	// fmt.Printf("Environment.Define: %v -> %v -> %v\n\n", name, value, e)
}

func (e *Environment) Get(name Token) (any, error) {
	if _, found := e.values[name.Lexeme]; found {
		// fmt.Printf("Environment.Get: %v -> %v -> %v\n\n", name, e.values[name.Lexeme], e)
		return e.values[name.Lexeme], nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, RuntimeError{name, fmt.Sprintf("Undefined variable '" + name.Lexeme + "'.")}
}

func (e *Environment) GetAt(distance int, name string) any {
	if value, found := e.ancestor(distance).values[name]; found {
		return value
	} else {
		panic(fmt.Sprintf("Environment.GetAt expected value for '%s' at distance '%d", name, distance))
	}
}

func (e *Environment) AssignAt(distance int, name Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) Assign(name Token, value any) error {
	if _, found := e.values[name.Lexeme]; found {
		e.values[name.Lexeme] = value
		// fmt.Printf("Environment.Assign: %v -> %v -> %v\n\n", name, value, e)
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}
