package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p AstPrinter) print(expr Expr) (any, error) {
	return expr.Accept(p)
}

func (p AstPrinter) VisitBinaryExpr(expr Binary) (any, error) {
	return p.Parenthesize(expr.Operator.Lexeme, []Expr{expr.Left, expr.Right})
}

func (p AstPrinter) VisitGroupingExpr(expr Grouping) (any, error) {
	return p.Parenthesize("group", []Expr{expr.Expression})
}

func (p AstPrinter) VisitLiteralExpr(expr Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return expr.Value, nil
}

func (p AstPrinter) VisitUnaryExpr(expr Unary) (any, error) {
	return p.Parenthesize(expr.Operator.Lexeme, []Expr{expr.Right})
}

func (p AstPrinter) Parenthesize(name string, exprs []Expr) (any, error) {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		accepted, err := expr.Accept(p)
		if err != nil {
			return nil, err
		}
		sb.WriteString(fmt.Sprintf("%v", accepted))
	}
	sb.WriteString(")")
	return sb.String(), nil
}
