package main

type StmtVisitor interface {
	VisitBlockStmt(stmt BlockStmt) (any, error)
	VisitExpressionStmt(stmt ExpressionStmt) (any, error)
	VisitPrintStmt(stmt PrintStmt) (any, error)
	VisitVarStmt(stmt VarStmt) (any, error)
}

type Stmt interface {
	AcceptStmt(visitor StmtVisitor) (any, error)
}

type BlockStmt struct {
	Statements []Stmt
}

func (expr BlockStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(expr)
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

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (expr VarStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarStmt(expr)
}
