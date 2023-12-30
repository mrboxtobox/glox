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
	if p.matchSingle(For) {
		return p.forStatement()
	}
	if p.matchSingle(If) {
		return p.ifStatement()
	}
	if p.matchSingle(Print) {
		return p.printStatement()
	}
	if p.matchSingle(Return) {
		return p.returnStatement()
	}
	if p.matchSingle(While) {
		return p.whileStatement()
	}
	if p.matchSingle(LeftBrace) {
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
	if _, err := p.consume(LeftParen, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}
	var initializer Stmt
	if p.matchSingle(Semicolon) {
		// Pass.
	} else if p.matchSingle(Var) {
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
	if !p.check(Semicolon) {
		var err error
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(Semicolon, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RightParen) {
		var err error
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(RightParen, "Expect ')' after for clauses."); err != nil {
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
	_, err := p.consume(LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if p.matchSingle(Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return IfStmt{condition, thenBranch, elseBranch}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}
	_, err := p.consume(RightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.matchSingle(Fun) {
		function, err := p.function("function")
		if err != nil {
			return nil, err
		}
		return function, nil
	}
	if p.matchSingle(Var) {
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

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.matchSingle(Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(Semicolon, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return VarStmt{name, initializer}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	if _, err := p.consume(LeftParen, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(RightParen, "Expect ')' after condition."); err != nil {
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
	if _, err = p.consume(Semicolon, "Expect ';' after value."); err != nil {
		return nil, err
	}
	return PrintStmt{value}, nil
}

func (p *Parser) returnStatement() (ReturnStmt, error) {
	keyword := p.previous()
	var value Expr
	if !p.check(Semicolon) {
		var err error
		value, err = p.expression()
		if err != nil {
			return ReturnStmt{}, err
		}
	}
	if _, err := p.consume(Semicolon, "Expect ';' after return value."); err != nil {
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

	if _, err = p.consume(Semicolon, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ExpressionStmt{expr}, nil
}

func (p *Parser) function(kind string) (FunctionStmt, error) {
	name, err := p.consume(Identifier, "Expect "+kind+" name.")
	if err != nil {
		return FunctionStmt{}, err
	}
	if _, err := p.consume(LeftParen, "Expect '(' after "+kind+" name."); err != nil {
		return FunctionStmt{}, err
	}
	var parameters []Token
	// Do-While loop.
	if !p.check(RightParen) {
		for {
			if len(parameters) >= 255 {
				PrintDetailedError(p.peek(), "Can't have more than 255 parameters.")
			}
			token, err := p.consume(Identifier, "Expect parameter name.")
			if err != nil {
				return FunctionStmt{}, err
			}
			parameters = append(parameters, token)
			if !p.matchSingle(Comma) {
				break
			}
		}
	}
	// fmt.Printf("%v\n\n", p.peek())
	// fmt.Printf("%v\n\n", parameters)
	// fmt.Printf("%v\n\n", p.tokens)
	if _, err := p.consume(RightParen, "Expect ')' after parameters."); err != nil {
		return FunctionStmt{}, err
	}
	if _, err := p.consume(LeftBrace, "Expect '{' before "+kind+" body."); err != nil {
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

	if p.matchSingle(Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if ve, ok := expr.(VariableExpr); ok {
			name := ve.Name
			return AssignExpr{name, value}, nil
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
	for p.matchSingle(Or) {
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
	for p.matchSingle(And) {
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
	for p.match([]TokenType{BangEqual, EqualEqual}) {
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
	for p.match([]TokenType{Greater, GreaterEqual, Less, LessEqual}) {
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
	for p.match([]TokenType{Minus, Plus}) {
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
	for p.match([]TokenType{Slash, Star}) {
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
	if p.match([]TokenType{Bang, Minus}) {
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
	if !p.check(RightParen) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, expression)
		for p.matchSingle(Comma) {
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
	paren, err := p.consume(RightParen, "Expect ')' after arguments.")
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
		if p.matchSingle(LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.matchSingle(False) {
		return LiteralExpr{false}, nil
	}
	if p.matchSingle(True) {
		return LiteralExpr{true}, nil
	}
	if p.matchSingle(Nil) {
		return LiteralExpr{nil}, nil
	}
	if p.match([]TokenType{Number, String}) {
		return LiteralExpr{p.previous().Literal}, nil
	}
	if p.matchSingle(Identifier) {
		return VariableExpr{p.previous()}, nil
	}
	if p.matchSingle(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RightParen, "Expect ')' after expression but got."); err != nil {
			return nil, err
		}
		return GroupingExpr{expr}, nil
	}

	// NOTE: This deviates from Ch. 6 error reporting since Go does not support
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
	// NOTE: This deviates from Ch. 6 error reporting since Go does not support
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
	return p.peek().TokenType == EOF
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
		if p.previous().TokenType == Semicolon {
			return
		}

		switch p.peek().TokenType {
		case Class:
			return
		case Fun:
			return
		case Var:
			return
		case For:
			return
		case If:
			return
		case While:
			return
		case Print:
			return
		case Return:
			return
		}
		p.advance()
	}
}

// NOTE: For debugging.
func (p *Parser) String() string {
	return fmt.Sprintf("*Parser{tokens=%v, position=%d}", p.tokens, p.current)
}
