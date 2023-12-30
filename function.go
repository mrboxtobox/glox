package main

type Function struct {
	declaration FunctionStmt
	// Not completely confident about this.
	closure *Environment
}

func NewFunction(declaration FunctionStmt, closure *Environment) Function {
	return Function{
		declaration: declaration,
		closure:     closure,
	}
}

// Manage name environments.
func (f Function) Call(interpreter Interpreter, arguments []any) (any, error) {
	// Use lexical scope at declaration.
	environment := NewEnvironmentFromEnclosing(f.closure)
	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}
	err := interpreter.executeBlock(f.declaration.Body, environment)
	if err == nil {
		return nil, nil
	}
	// `FunctionReturn` is technically not an error.
	if fr, ok := err.(FunctionReturn); ok {
		return fr.Value, nil
	}
	return nil, err
}

func (f Function) Arity() int {
	return len(f.declaration.Params)
}

func (f Function) ToString() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
