package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Clock struct{}

func (Clock) Arity() int {
	return 0
}

func (Clock) Call(interpreter Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli()), nil
}

func (Clock) String() string {
	return "<native fn: clock>"
}

type Interpreter struct {
	environment *Environment
	// TODO: Figure out if we need globals.
	globals *Environment
	// string is the ptr of Expr
	locals map[string]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment()
	environment := globals
	globals.Define("clock", Clock{})
	return &Interpreter{
		environment: environment,
		globals:     environment,
		// Each expression node is its own object. No need for a nested
		// tree.
		locals: map[string]int{},
	}
}

func (i Interpreter) Interpret(statements []Stmt) error {
	for _, statement := range statements {
		if _, err := i.execute(statement); err != nil {
			switch typedErr := err.(type) {
			case FunctionReturn:
				continue
			case RuntimeError:
				println(fmt.Sprintf("[line %d] Runtime error: %s\n", typedErr.Token.Line, typedErr.Message))
				return typedErr
			default:
				return typedErr
			}
		}
	}
	return nil
}

func (i Interpreter) execute(statement Stmt) (any, error) {
	// println("Executing -> ")
	// fmt.Printf("  %T -> %v\n", statement, statement)
	return statement.AcceptStmt(i)
}

// TODO: Figure out if this hash is the correct way.
func (i Interpreter) Resolve(expr Expr, depth int) {
	i.locals[fmt.Sprintf("%v", expr)] = depth
}

func (i Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	// TODO: Figure out how this assignment will work.
	previous := i.environment
	// Use defer to ensure the environment is restored even in case of an early return.
	defer func() {
		i.environment = previous
	}()
	i.environment = environment
	for _, statement := range statements {
		if _, err := i.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (i Interpreter) VisitBlockStmt(stmt BlockStmt) (any, error) {
	if err := i.executeBlock(stmt.Statements, NewEnvironmentFromEnclosing(i.environment)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i Interpreter) VisitClassStmt(stmt ClassStmt) (any, error) {
	var superclass any
	if stmt.Superclass != (VariableExpr{}) {
		var err error
		superclass, err = i.evaluate(stmt.Superclass)
		if err != nil {
			return nil, err
		}

		if _, ok := superclass.(*Class); !ok {
			err := LogAndReturnError(stmt.Superclass.Name, "Superclass must be a class.")
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)
	if stmt.Superclass != (VariableExpr{}) {
		i.environment = NewEnvironmentFromEnclosing(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := map[string]Function{}
	for _, method := range stmt.Methods {
		isInitializer := method.Name.Lexeme == "init"
		function := NewFunction(method, i.environment, isInitializer)
		methods[method.Name.Lexeme] = function
	}
	var class *Class
	if superclass == nil {
		class = NewClass(stmt.Name.Lexeme, nil, methods)
	} else {
		class = NewClass(stmt.Name.Lexeme, superclass.(*Class), methods)
	}
	if superclass != nil {
		i.environment = i.environment.enclosing
	}
	if err := i.environment.Assign(stmt.Name, class); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i Interpreter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i Interpreter) VisitLogicalExpr(expr LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == OrToken {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

func (i Interpreter) VisitSetExpr(expr SetExpr) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}
	instance, ok := object.(*Instance)
	if !ok {
		err := LogAndReturnError(expr.Name, "Only instances have fields.")
		return nil, err
	}
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	instance.Set(expr.Name, value)
	return value, nil
}

func (i Interpreter) VisitSuperExpr(expr SuperExpr) (any, error) {
	distance, found := i.locals[fmt.Sprintf("%v", expr)]
	if !found {
		return nil, fmt.Errorf("Expected %v to be found in locals\n", expr)
	}
	superclass := i.environment.GetAt(distance, "super").(*Class)
	// The environment where “this” is bound is always right inside the environment where we store “super”.
	this := i.environment.GetAt(distance-1, "this").(*Instance)
	method, found := superclass.FindMethod(expr.Method.Lexeme)
	if !found {
		err := LogAndReturnError(expr.Method, "Undefined property '"+expr.Keyword.Lexeme+"'.")
		return nil, err
	}
	return method.Bind(this), nil
}

func (i Interpreter) VisitThisExpr(expr ThisExpr) (any, error) {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i Interpreter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.TokenType {
	case MinusToken:
		if err := checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -(right.(float64)), nil
	case BangToken:
		return !isTruthy(right), nil
	}
	// Unreachable.
	return nil, nil
}

func (i Interpreter) VisitBinaryExpr(expr BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case GreaterToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqualToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LessToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqualToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	// NOTE: The two cases below differ from the Ch. 7 Java implementation.
	// IMPORTANT: NaN != NaN according to the IEEE spec.
	case BangEqualToken:
		return left != right, nil
	case EqualEqualToken:
		return left == right, nil
	case MinusToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case PlusToken:
		leftKind := reflect.TypeOf(left).Kind()
		rightKind := reflect.TypeOf(right).Kind()
		if leftKind == reflect.Float64 && rightKind == reflect.Float64 {
			return left.(float64) + right.(float64), nil
		}
		if leftKind == reflect.String && rightKind == reflect.String {
			return left.(string) + right.(string), nil
		}
		return nil, RuntimeError{expr.Operator, fmt.Sprintf("Operands (%v, %v) must be two numbers or two strings", left, right)}
	case SlashToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case StarToken:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}
	// Unreachable.
	return nil, nil
}

func (i Interpreter) VisitCallExpr(expr CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	var arguments []any
	for _, argument := range expr.Arguments {
		value, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, value)
	}
	function, ok := callee.(Callable)
	if !ok {
		err := LogAndReturnError(expr.Paren, "Can only call functions and classes.")
		return nil, err
	}
	if len(arguments) != function.Arity() {
		err := LogAndReturnError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", +function.Arity(), len(arguments)))
		return nil, err
	}
	return function.Call(i, arguments)
}

func (i Interpreter) VisitGetExpr(expr GetExpr) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}
	if instance, ok := object.(*Instance); ok {
		return instance.Get(expr.Name)
	}
	err = LogAndReturnError(expr.Name, "Only instances have properties.")
	return nil, err
}

func (i Interpreter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitVariableExpr(expr VariableExpr) (any, error) {
	return i.lookUpVariable(expr.Name, expr)
}

func (i Interpreter) lookUpVariable(name Token, expr Expr) (any, error) {
	if distance, found := i.locals[fmt.Sprintf("%v", expr)]; found {
		return i.environment.GetAt(distance, name.Lexeme), nil
	} else {
		return i.globals.Get(name)
	}
}

func (i Interpreter) evaluate(expr Expr) (any, error) {
	return expr.AcceptExpr(i)
}

func (i Interpreter) VisitExpressionStmt(stmt ExpressionStmt) (any, error) {
	// There's no result from a statement. So, just evaluate and ignore the result.
	if _, err := i.evaluate(stmt.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i Interpreter) VisitFunctionStmt(stmt FunctionStmt) (any, error) {
	// Choose the encironment that is active when the function is declared, not
	// called. Lexical scope surrounding the function declaration.
	// TODO: Figure out why.
	// Lexical scope instead of globals.
	function := NewFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i Interpreter) VisitIfStmt(stmt IfStmt) (any, error) {
	value, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(value) {
		if _, err := i.execute(stmt.ThenBranch); err != nil {
			return nil, err
		}
	} else if stmt.ElseBranch != nil {
		if _, err := i.execute(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i Interpreter) VisitPrintStmt(stmt PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	println(stringify(value))
	return nil, nil
}

func (i Interpreter) VisitReturnStmt(stmt ReturnStmt) (any, error) {
	var value any
	if stmt.Value != nil {
		var err error
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	// Go does not allow exceptions so return an error. Errors are propagated all
	// the way up the stack(?).
	// TODO: Verify if this is true.
	return nil, FunctionReturn{value}
}

func (i Interpreter) VisitVarStmt(stmt VarStmt) (any, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

// TODO: Figure out if we need to break early.
func (i Interpreter) VisitWhileStmt(stmt WhileStmt) (any, error) {
	for {
		value, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !isTruthy(value) {
			return nil, nil
		}
		if _, err = i.execute(stmt.Body); err != nil {
			return nil, err
		}
	}
}

func (i Interpreter) VisitAssignExpr(expr AssignExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if distance, found := i.locals[fmt.Sprintf("%v", expr)]; found {
		i.environment.AssignAt(distance, expr.Name, value)
		return value, nil
	}

	if err = i.globals.Assign(expr.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if _, ok := value.(bool); ok {
		return value.(bool)
	}
	return true
}

func checkNumberOperand(operator Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	// TODO: Need to figure out how to throw errors properly.
	return RuntimeError{operator, fmt.Sprintf("Operand '%v' must be a number but is '%T'", operand, operand)}
}

func checkNumberOperands(operator Token, left any, right any) error {
	leftKind := reflect.TypeOf(left).Kind()
	rightKind := reflect.TypeOf(right).Kind()
	if leftKind == reflect.Float64 && rightKind == reflect.Float64 {
		return nil
	}
	return RuntimeError{operator, fmt.Sprintf("Operands (%v, %v) must be numbers but are (%T, %T).", left, right, left, right)}
}

// NOTE: Update this for any custom type that we want .
func stringify(object any) string {
	if object == nil {
		return "nil"
	}

	switch object := object.(type) {
	case float64:
		return strconv.FormatFloat(object, 'f', -1, 64)
	}
	if stringer, ok := object.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", object)
}
