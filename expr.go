package main

type ExprVisitor interface {
	VisitBinaryExpr(expr BinaryExpr) (any, error)
	VisitGroupingExpr(expr GroupingExpr) (any, error)
	VisitLiteralExpr(expr LiteralExpr) (any, error)
	VisitUnaryExpr(expr UnaryExpr) (any, error)
}

type Expr interface {
	AcceptExpr(visitor ExprVisitor) (any, error)
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr BinaryExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(expr)
}

type GroupingExpr struct {
	Expression Expr
}

func (expr GroupingExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(expr)
}

type LiteralExpr struct {
	Value interface{}
}

func (expr LiteralExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr UnaryExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(expr)
}
