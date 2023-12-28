package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Interpreter struct {
	environment Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		if _, err := i.execute(statement); err != nil {
			PrintRuntimeError(err.(RuntimeError))
			break
		}
	}
}

func (i Interpreter) execute(statement Stmt) (any, error) {
	return statement.AcceptStmt(i)
}

func (i Interpreter) executeBlock(statements []Stmt, environment Environment) error {
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
	if err := i.executeBlock(stmt.Statements, NewEnvironmentFromEnclosing(&i.environment)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i Interpreter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i Interpreter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.TokenType {
	case Minus:
		if err := checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -(right.(float64)), nil
	case Bang:
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
	case Greater:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case Less:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	// NOTE: The two cases below differ from the Ch. 7 Java implementation.
	// IMPORTANT: NaN != NaN according to the IEEE spec.
	case BangEqual:
		return left == right, nil
	case EqualEqual:
		return left != right, nil
	case Minus:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case Plus:
		leftKind := reflect.TypeOf(left).Kind()
		rightKind := reflect.TypeOf(right).Kind()
		if leftKind == reflect.Float64 && rightKind == reflect.Float64 {
			return left.(float64) + right.(float64), nil
		}
		if leftKind == reflect.String && rightKind == reflect.String {
			return left.(string) + right.(string), nil
		}
		return nil, RuntimeError{expr.Operator, fmt.Sprintf("Operands (%q, %q) must be two numbers or two strings", left, right)}
	case Slash:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case Star:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}
	// Unreachable.
	return nil, nil
}

func (i Interpreter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitVariableExpr(expr VariableExpr) (any, error) {
	return i.environment.Get(expr.Name)
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

func (i Interpreter) VisitPrintStmt(stmt PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	println(stringify(value))
	return nil, nil
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

func (i Interpreter) VisitAssignExpr(expr AssignExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.environment.Assign(expr.Name, value)
	if err != nil {
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
	return RuntimeError{operator, fmt.Sprintf("Operands (%v, %v) must be numbers but are (%T, %T)", left, right, left, right)}
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}

	switch object.(type) {
	case float64:
		text := fmt.Sprintf("%f", object)
		text, _ = strings.CutSuffix(text, ".000000")
		return text
	}
	return fmt.Sprintf("%v", object)
}
