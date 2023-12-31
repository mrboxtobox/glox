package main

type Function struct {
	declaration FunctionStmt
	// Not completely confident about this.
	closure       *Environment
	isInitializer bool
}

func NewFunction(declaration FunctionStmt, closure *Environment, isInitializer bool) Function {
	return Function{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
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
		if f.isInitializer {
			return f.closure.GetAt(0, "this"), nil
		}
		return fr.Value, nil
	}
	if f.isInitializer {
		return f.closure.GetAt(0, "this"), nil
	}
	return nil, err
}

func (f Function) Arity() int {
	return len(f.declaration.Params)
}

func (f Function) Bind(instance *Instance) Function {
	environment := NewEnvironmentFromEnclosing(f.closure)
	environment.Define("this", instance)
	return NewFunction(f.declaration, environment, f.isInitializer)
}

func (f Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
