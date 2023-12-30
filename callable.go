package main

type Callable interface {
	Arity() int
	Call(interpreter Interpreter, arguments []any) (any, error)
	ToString() string
}
