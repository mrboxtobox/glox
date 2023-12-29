package main

type StmtVisitor interface {
	VisitBlockStmt(stmt BlockStmt) (any, error)
	VisitExpressionStmt(stmt ExpressionStmt) (any, error)
	VisitIfStmt(stmt IfStmt) (any, error)
	VisitPrintStmt(stmt PrintStmt) (any, error)
	VisitVarStmt(stmt VarStmt) (any, error)
	VisitWhileStmt(stmt WhileStmt) (any, error)
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

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (expr IfStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitIfStmt(expr)
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

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (expr WhileStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(expr)
}
