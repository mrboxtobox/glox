package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p AstPrinter) print(expr Expr) any {
	return expr.Accept(p)
}

func (p AstPrinter) VisitBinaryExpr(expr Binary) any {
	return p.Parenthesize(expr.Operator.Lexeme, []Expr{expr.Left, expr.Right})
}

func (p AstPrinter) VisitGroupingExpr(expr Grouping) any {
	return p.Parenthesize("group", []Expr{expr.Expression})
}

func (p AstPrinter) VisitLiteralExpr(expr Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value
}

func (p AstPrinter) VisitUnaryExpr(expr Unary) any {
	return p.Parenthesize(expr.Operator.Lexeme, []Expr{expr.Right})
}

func (p AstPrinter) Parenthesize(name string, exprs []Expr) any {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", expr.Accept(p)))
	}
	sb.WriteString(")")
	return sb.String()
}
