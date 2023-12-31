package main

import "fmt"

type Callable interface {
	Arity() int
	Call(interpreter Interpreter, arguments []any) (any, error)
	fmt.Stringer
}
