package main

import (
	"strings"
)

type AstPrinter struct{}

func (p AstPrinter) print(expr Expr) {
	return expr.accept(p)
}

func (p AstPrinter) visitBinaryExpr(expr Binary) string {
	return p.parenthesize(expr.operator.lexeme, []Expr{expr.left, expr.right})
}

func (p AstPrinter) visitGroupingExpr(expr Grouping) string {
	return p.parenthesize("group", expr.expression)
}

func (p AstPrinter) visitLiteralExpr(expr Literal) string {
	if expr.value == nil {
		return "nil"
	}
	return expr.value
}

func (p AstPrinter) visitUnaryExpr(expr Unary) string {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

func (p AstPrinter) parenthesize(name string, exprs []Expr) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString("(")
		sb.WriteString(expr.accept(p))
	}
	sb.WriteString(")")
	return sb.String()
}
