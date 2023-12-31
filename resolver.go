package main

import "fmt"

type FunctionType int
type ClassType int

const (
	NoneFunction FunctionType = iota
	InFunction
	InitializerFunction
	Method
)

const (
	NoneClass ClassType = iota
	InClass
	SubClass
)

type Resolver struct {
	interpreter     Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver(interpreter Interpreter) Resolver {
	return Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NoneFunction,
		currentClass:    NoneClass,
	}
}

func (r *Resolver) VisitAssignExpr(expr AssignExpr) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt VarStmt) (any, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil, nil
}

// VisitBinaryExpr implements ExprVisitor.
func (r *Resolver) VisitBinaryExpr(expr BinaryExpr) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr CallExpr) (any, error) {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr GetExpr) (any, error) {
	r.resolveExpr(expr.Object)
	// Properties are resolved dynamically.
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	r.resolveExpr(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr LogicalExpr) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr SetExpr) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr SuperExpr) (any, error) {
	if r.currentClass == NoneClass {
		PrintDetailedError(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != SubClass {
		PrintDetailedError(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr ThisExpr) (any, error) {
	if r.currentClass == NoneClass {
		PrintDetailedError(expr.Keyword, "Can't use 'this' outside of a class.")
		return nil, nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

// VisitUnaryExpr implements ExprVisitor.
func (r *Resolver) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

// VisitVariableExpr implements ExprVisitor.
func (r *Resolver) VisitVariableExpr(expr VariableExpr) (any, error) {
	if len(r.scopes) > 0 && !r.scopes[len(r.scopes)-1][expr.Name.Lexeme] {
		// If the varaible exists in the scope and it's value is false, print an
		// error to indicate it shouldn't be used.
		PrintDetailedError(expr.Name, "Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) resolveAll(statements []Stmt) (any, error) {
	for _, statement := range statements {
		if err := r.resolveStmt(statement); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitBlockStmt(stmt BlockStmt) (any, error) {
	r.beginScope()
	if _, err := r.resolveAll(stmt.Statements); err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt ClassStmt) (any, error) {
	enclosingClass := r.currentClass
	r.currentClass = InClass
	r.declare(stmt.Name)
	r.define(stmt.Name)
	if stmt.Superclass != (VariableExpr{}) && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		PrintDetailedError(stmt.Superclass.Name, "A class can't inherit from itself.")
	}

	if stmt.Superclass != (VariableExpr{}) {
		r.currentClass = SubClass
		r.resolveExpr(stmt.Superclass)
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}
	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true
	for _, method := range stmt.Methods {
		declaration := Method
		if method.Name.Lexeme == "init" {
			declaration = InitializerFunction
		}
		if err := r.resolveFunction(method, declaration); err != nil {
			return nil, err
		}
	}
	r.endScope()
	if stmt.Superclass != (VariableExpr{}) {
		r.endScope()
	}
	r.currentClass = enclosingClass
	return nil, nil
}

// VisitExpressionStmt implements StmtVisitor.
func (r *Resolver) VisitExpressionStmt(stmt ExpressionStmt) (any, error) {
	r.resolveExpr(stmt.Expression)
	return nil, nil
}

// VisitFunctionStmt implements StmtVisitor.
func (r *Resolver) VisitFunctionStmt(stmt FunctionStmt) (any, error) {
	// Declare and define first.
	// This lets a function recursively refer to itself inside its own body.
	r.declare(stmt.Name)
	r.define(stmt.Name)
	if err := r.resolveFunction(stmt, InFunction); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt IfStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt PrintStmt) (any, error) {
	r.resolveExpr(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt ReturnStmt) (any, error) {
	if r.currentFunction == NoneFunction {
		PrintDetailedError(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		if r.currentFunction == InitializerFunction {
			PrintDetailedError(stmt.Keyword, "Can't return a value from an initializer.")
		}
		r.resolveExpr(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt WhileStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	if err := r.resolveStmt(stmt.Body); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) resolveStmt(stmt Stmt) error {
	if _, err := stmt.AcceptStmt(r); err != nil {
		// TODO: Handle error properly.
		PrintResolverError(err)
		return err
	}
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1] // Peek
	if _, found := scope[name.Lexeme]; found {
		PrintDetailedError(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	n := len(r.scopes)
	for i := n - 1; i >= 0; i-- {
		if _, found := r.scopes[i][name.Lexeme]; found {
			r.interpreter.Resolve(expr, n-1-i) // n-1-i = Number of hops
			return
		}
	}
	// Leave unresolved and assume it's global.

}

func (r *Resolver) resolveExpr(expr Expr) {
	if _, err := expr.AcceptExpr(r); err != nil {
		// TODO: Handle error properly.
		PrintResolverError(err)
	}
}

func (r *Resolver) resolveFunction(function FunctionStmt, typ FunctionType) error {
	// Stash previous value of the field in local variable first.
	enclosingFunction := r.currentFunction
	r.currentFunction = typ
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	if _, err := r.resolveAll(function.Body); err != nil {
		return err
	}
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func PrintResolverError(err error) {
	fmt.Printf("Error during variable resolution: %v\n", err)
}
