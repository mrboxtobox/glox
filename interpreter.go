package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Interpreter struct{}

func (p Interpreter) Interpret(expr Expr) {
	value, err := p.evaluate(expr)
	if err != nil {
		// TODO: There should be a better to way to do the cast.
		PrintRuntimeError(err.(RuntimeError))
	}
	fmt.Printf("%v\n", stringify(value))
}

func (p Interpreter) VisitLiteralExpr(expr Literal) (any, error) {
	return expr.Value, nil
}

func (p Interpreter) VisitUnaryExpr(expr Unary) (any, error) {
	right, err := p.evaluate(expr.Right)
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

func (p Interpreter) VisitBinaryExpr(expr Binary) (any, error) {
	left, err := p.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := p.evaluate(expr.Right)
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

func (p Interpreter) VisitGroupingExpr(expr Grouping) (any, error) {
	return p.evaluate(expr.Expression)
}

func (p Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(p)
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
	return RuntimeError{operator, fmt.Sprintf("Operands (%q, %q) must be numbers but are (%T, %T)", left, right, left, right)}
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
	return fmt.Sprintf("%s", object)
}
