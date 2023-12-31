package main

type StmtVisitor interface {
	VisitBlockStmt(stmt BlockStmt) (any, error)
	VisitClassStmt(stmt ClassStmt) (any, error)
	VisitExpressionStmt(stmt ExpressionStmt) (any, error)
	VisitFunctionStmt(stmt FunctionStmt) (any, error)
	VisitIfStmt(stmt IfStmt) (any, error)
	VisitPrintStmt(stmt PrintStmt) (any, error)
	VisitVarStmt(stmt VarStmt) (any, error)
	VisitReturnStmt(stmt ReturnStmt) (any, error)
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

type ClassStmt struct {
	Name       Token
	Superclass VariableExpr
	Methods    []FunctionStmt
}

func (expr ClassStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitClassStmt(expr)
}

type ExpressionStmt struct {
	Expression Expr
}

func (expr ExpressionStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(expr)
}

type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (expr FunctionStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitFunctionStmt(expr)
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

type ReturnStmt struct {
	Keyword Token
	Value   Expr
}

func (expr ReturnStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitReturnStmt(expr)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (expr WhileStmt) AcceptStmt(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(expr)
}
