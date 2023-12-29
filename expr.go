package main

type ExprVisitor interface {
	VisitAssignExpr(expr AssignExpr) (any, error)
	VisitBinaryExpr(expr BinaryExpr) (any, error)
	VisitGroupingExpr(expr GroupingExpr) (any, error)
	VisitLiteralExpr(expr LiteralExpr) (any, error)
	VisitLogicalExpr(expr LogicalExpr) (any, error)
	VisitUnaryExpr(expr UnaryExpr) (any, error)
	VisitVariableExpr(expr VariableExpr) (any, error)
}

type Expr interface {
	AcceptExpr(visitor ExprVisitor) (any, error)
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (expr AssignExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssignExpr(expr)
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

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr LogicalExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr UnaryExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(expr)
}

type VariableExpr struct {
	Name Token
}

func (expr VariableExpr) AcceptExpr(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(expr)
}
