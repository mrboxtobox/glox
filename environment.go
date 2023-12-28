package main

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() Environment {
	return Environment{
		values: map[string]any{},
	}
}

func NewEnvironmentFromEnclosing(enclosing *Environment) Environment {
	return Environment{
		enclosing: enclosing,
		values:    map[string]any{},
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	if _, found := e.values[name.Lexeme]; found {
		return e.values[name.Lexeme], nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, RuntimeError{name, fmt.Sprintf("Undefined variable '" + name.Lexeme + "'.")}
}

func (e *Environment) Assign(name Token, value any) error {
	if _, found := e.values[name.Lexeme]; found {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}
