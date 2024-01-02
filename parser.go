package main

import "fmt"

type ParseError struct{}

func (e ParseError) Error() string {
	return "parse error"
}

type Parser struct {
	// Tokens to parse.
	tokens []Token
	// Points to the next token to be parsed.
	current int
}

// Use a pointer receiver to ensure that methods can modify the values.
// See https://go.dev/tour/methods/8.
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}
	return statements, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.matchSingle(ForToken) {
		return p.forStatement()
	}
	if p.matchSingle(IfToken) {
		return p.ifStatement()
	}
	if p.matchSingle(PrintToken) {
		return p.printStatement()
	}
	if p.matchSingle(ReturnToken) {
		return p.returnStatement()
	}
	if p.matchSingle(WhileToken) {
		return p.whileStatement()
	}
	if p.matchSingle(LeftBraceToken) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return BlockStmt{statements}, nil
	}
	return p.expressionStatement()
}

// Desugaring by transforming a for loop into a while loop.
func (p *Parser) forStatement() (Stmt, error) {
	if _, err := p.consume(LeftParenToken, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}
	var initializer Stmt
	if p.matchSingle(SemicolonToken) {
		// Pass.
	} else if p.matchSingle(VarToken) {
		var err error
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	// Check if a semicolon immediately follows the last clause.
	if !p.check(SemicolonToken) {
		var err error
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(SemicolonToken, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RightParenToken) {
		var err error
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(RightParenToken, "Expect ')' after for clauses."); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	// Build up the components of the for loop using while semantics.
	// TODO: check if the nil value is Expr.
	if increment != nil {
		body = BlockStmt{[]Stmt{body, ExpressionStmt{increment}}}
	}
	if condition == nil {
		condition = LiteralExpr{true}
	}
	body = WhileStmt{condition, body}
	if initializer != nil {
		body = BlockStmt{[]Stmt{initializer, body}}
	}
	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LeftParenToken, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParenToken, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if p.matchSingle(ElseToken) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return IfStmt{condition, thenBranch, elseBranch}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt
	for !p.check(RightBraceToken) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}
	_, err := p.consume(RightBraceToken, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.matchSingle(ClassToken) {
		class, err := p.classDeclaration()
		if err != nil {
			return nil, err
		}
		return class, nil
	}
	if p.matchSingle(FunToken) {
		function, err := p.function("function")
		if err != nil {
			return nil, err
		}
		return function, nil
	}
	if p.matchSingle(VarToken) {
		declaration, err := p.varDeclaration()
		if err == nil {
			return declaration, nil
		}
		// TODO: Handle non ParseErrors.
		if _, ok := err.(ParseError); ok {
			p.synchronize()
			return nil, err
		} else {
			panic(err)
		}
	}

	statement, err := p.statement()
	if err == nil {
		return statement, err
	}
	if _, ok := err.(ParseError); ok {
		p.synchronize()
		return nil, err
	} else {
		panic(err)
	}
}

func (p *Parser) classDeclaration() (Stmt, error) {
	name, err := p.consume(IdentifierToken, "Expect class name.")
	if err != nil {
		return nil, err
	}
	var superclass VariableExpr
	if p.matchSingle(LessToken) {
		if _, err := p.consume(IdentifierToken, "Expect superclass name."); err != nil {
			return nil, err
		}
		superclass = VariableExpr{p.previous()}
	}
	if _, err = p.consume(LeftBraceToken, "Expect '{' before class body."); err != nil {
		return nil, err
	}
	var methods []FunctionStmt
	for !p.check(RightBraceToken) && !p.isAtEnd() {
		method, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	if _, err := p.consume(RightBraceToken, "Expect '}' after class body."); err != nil {
		return nil, err
	}
	return ClassStmt{name, superclass, methods}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IdentifierToken, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.matchSingle(EqualToken) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(SemicolonToken, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return VarStmt{name, initializer}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	if _, err := p.consume(LeftParenToken, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(RightParenToken, "Expect ')' after condition."); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return WhileStmt{condition, body}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(SemicolonToken, "Expect ';' after value."); err != nil {
		return nil, err
	}
	return PrintStmt{value}, nil
}

func (p *Parser) returnStatement() (ReturnStmt, error) {
	keyword := p.previous()
	var value Expr
	if !p.check(SemicolonToken) {
		var err error
		value, err = p.expression()
		if err != nil {
			return ReturnStmt{}, err
		}
	}
	if _, err := p.consume(SemicolonToken, "Expect ';' after return value."); err != nil {
		return ReturnStmt{}, err
	}
	return ReturnStmt{keyword, value}, nil

}

// TODO: Change the return types to be concrete.
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(SemicolonToken, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ExpressionStmt{expr}, nil
}

func (p *Parser) function(kind string) (FunctionStmt, error) {
	name, err := p.consume(IdentifierToken, "Expect "+kind+" name.")
	if err != nil {
		return FunctionStmt{}, err
	}
	if _, err := p.consume(LeftParenToken, "Expect '(' after "+kind+" name."); err != nil {
		return FunctionStmt{}, err
	}
	var parameters []Token
	// Do-WhileToken loop.
	if !p.check(RightParenToken) {
		for {
			if len(parameters) >= 255 {
				PrintDetailedError(p.peek(), "Can't have more than 255 parameters.")
			}
			token, err := p.consume(IdentifierToken, "Expect parameter name.")
			if err != nil {
				return FunctionStmt{}, err
			}
			parameters = append(parameters, token)
			if !p.matchSingle(CommaToken) {
				break
			}
		}
	}
	// fmt.Printf("%v\n\n", p.peek())
	// fmt.Printf("%v\n\n", parameters)
	// fmt.Printf("%v\n\n", p.tokens)
	if _, err := p.consume(RightParenToken, "Expect ')' after parameters."); err != nil {
		return FunctionStmt{}, err
	}
	if _, err := p.consume(LeftBraceToken, "Expect '{' before "+kind+" body."); err != nil {
		return FunctionStmt{}, err
	}
	body, err := p.block()
	if err != nil {
		return FunctionStmt{}, err
	}
	return FunctionStmt{name, parameters, body}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// Assignment is right-associative.
// We can do this since every valid assignment target is a valid expression.
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.matchSingle(EqualToken) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if ve, ok := expr.(VariableExpr); ok {
			name := ve.Name
			return AssignExpr{name, value}, nil
		} else if get, ok := expr.(GetExpr); ok {
			return SetExpr{get.Object, get.Name, value}, nil
		}
		// We don't throw an error because the parser is not in a bad state.
		PrintDetailedError(equals, "Invalid assignment target.")
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.matchSingle(OrToken) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = LogicalExpr{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.matchSingle(AndToken) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = LogicalExpr{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match([]TokenType{BangEqualToken, EqualEqualToken}) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match([]TokenType{GreaterToken, GreaterEqualToken, LessToken, LessEqualToken}) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match([]TokenType{MinusToken, PlusToken}) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match([]TokenType{SlashToken, StarToken}) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match([]TokenType{BangToken, MinusToken}) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{operator, right}, nil
	}
	return p.call()
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr
	if !p.check(RightParenToken) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, expression)
		for p.matchSingle(CommaToken) {
			// Go doesn't seem to have a limit. So, use 255 (Java's limit).
			// Only report an error but don't throw since the Parser is in a
			// valid state.
			if len(arguments) > 255 {
				PrintDetailedError(p.peek(), "Can't have more than 255 arguments.")
			}
			expression, err = p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expression)
		}
	}
	paren, err := p.consume(RightParenToken, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return CallExpr{callee, paren, arguments}, nil
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.matchSingle(LeftParenToken) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.matchSingle(DotToken) {
			name, err := p.consume(IdentifierToken, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = GetExpr{expr, name}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.matchSingle(FalseToken) {
		return LiteralExpr{false}, nil
	}
	if p.matchSingle(TrueToken) {
		return LiteralExpr{true}, nil
	}
	if p.matchSingle(NilToken) {
		return LiteralExpr{nil}, nil
	}
	if p.match([]TokenType{NumberToken, StringToken}) {
		return LiteralExpr{p.previous().Literal}, nil
	}
	if p.matchSingle(SuperToken) {
		keyword := p.previous()
		if _, err := p.consume(DotToken, "Expect '.' after 'super'."); err != nil {
			return nil, err
		}
		method, err := p.consume(IdentifierToken, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return SuperExpr{keyword, method}, nil
	}
	if p.matchSingle(ThisToken) {
		return ThisExpr{p.previous()}, nil
	}
	if p.matchSingle(IdentifierToken) {
		return VariableExpr{p.previous()}, nil
	}
	if p.matchSingle(LeftParenToken) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RightParenToken, "Expect ')' after expression but got."); err != nil {
			return nil, err
		}
		return GroupingExpr{expr}, nil
	}

	// NOTE: ThisToken deviates from Ch. 6 error reporting since Go does not support
	// throwing errors.
	PrintDetailedError(p.peek(), fmt.Sprintf("Expect expression but got '%v'.", p.peek()))
	return nil, ParseError{}
}

func (p *Parser) match(types []TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) matchSingle(tokenType TokenType) bool {
	if p.check(tokenType) {
		p.advance()
		return true
	}
	return false
}

// TODO: Use the value of `consume`.
func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	// NOTE: ThisToken deviates from Ch. 6 error reporting since Go does not support
	// throwing errors.
	PrintDetailedError(p.peek(), message)
	return Token{}, ParseError{}
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOFToken
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// NOTE: error() should have been here.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SemicolonToken {
			return
		}

		switch p.peek().TokenType {
		case ClassToken:
			return
		case FunToken:
			return
		case VarToken:
			return
		case ForToken:
			return
		case IfToken:
			return
		case WhileToken:
			return
		case PrintToken:
			return
		case ReturnToken:
			return
		}
		p.advance()
	}
}

// NOTE: ForToken debugging.
func (p *Parser) String() string {
	return fmt.Sprintf("*Parser{tokens=%v, position=%d}", p.tokens, p.current)
}
