package main

type Visitor interface {
	visitBinaryExpr(expr Binary) any
	visitGroupingExpr(expr Grouping) any
	visitLiteralExpr(expr Literal) any
	visitUnaryExpr(expr Unary) any
}

type Expr interface {
	accept(visitor Visitor) any
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (expr Binary) accept(visitor Visitor) any {
	return visitor.visitBinaryExpr(expr)
}

type Grouping struct {
	expression Expr
}

func (expr Grouping) accept(visitor Visitor) any {
	return visitor.visitGroupingExpr(expr)
}

type Literal struct {
	value interface{}
}

func (expr Literal) accept(visitor Visitor) any {
	return visitor.visitLiteralExpr(expr)
}

type Unary struct {
	operator Token
	right    Expr
}

func (expr Unary) accept(visitor Visitor) any {
	return visitor.visitUnaryExpr(expr)
}
