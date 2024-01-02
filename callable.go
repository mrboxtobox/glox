package main

import "fmt"

// Something that can be 'called'. Either a function or method.
type Callable interface {
	// The number of arguments expected by the Callable.
	Arity() int
	// Executes the Callable.
	Call(interpreter Interpreter, arguments []any) (any, error)
	fmt.Stringer
}
