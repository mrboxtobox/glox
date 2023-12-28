package main

type Visitor interface {
	VisitBinaryExpr(expr Binary) (any, error)
	VisitGroupingExpr(expr Grouping) (any, error)
	VisitLiteralExpr(expr Literal) (any, error)
	VisitUnaryExpr(expr Unary) (any, error)
}

type Expr interface {
	Accept(visitor Visitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr Binary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinaryExpr(expr)
}

type Grouping struct {
	Expression Expr
}

func (expr Grouping) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGroupingExpr(expr)
}

type Literal struct {
	Value interface{}
}

func (expr Literal) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLiteralExpr(expr)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (expr Unary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnaryExpr(expr)
}
