package main

type Visitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
}

type Expr interface {
	Accept(visitor Visitor) any
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(expr)
}

type Grouping struct {
	Expression Expr
}

func (expr Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(expr)
}

type Literal struct {
	Value interface{}
}

func (expr Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(expr)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (expr Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(expr)
}
