package main

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt) (any, error)
	VisitPrintStmt(stmt PrintStmt) (any, error)
}

type Stmt interface {
	AcceptStmt(visitor StmtVisitor) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (expr ExpressionStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(expr)
}

type PrintStmt struct {
	Expression Expr
}

func (expr PrintStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(expr)
}
