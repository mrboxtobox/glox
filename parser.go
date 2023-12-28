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
	if p.match([]TokenType{Print}) {
		return p.printStatement()
	}
	if p.match([]TokenType{LeftBrace}) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return BlockStmt{statements}, nil
	}
	return p.expressionStatement()
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
	if p.match([]TokenType{Var}) {
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
	if p.match([]TokenType{Equal}) {
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

// Assignment is right-associative.
// We can do this since every valid assignment target is a valid expression.
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match([]TokenType{Equal}) {
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

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
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
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match([]TokenType{False}) {
		return LiteralExpr{false}, nil
	}
	if p.match([]TokenType{True}) {
		return LiteralExpr{true}, nil
	}
	if p.match([]TokenType{Nil}) {
		return LiteralExpr{nil}, nil
	}
	if p.match([]TokenType{Number, String}) {
		return LiteralExpr{p.previous().Literal}, nil
	}
	if p.match([]TokenType{Identifier}) {
		return VariableExpr{p.previous()}, nil
	}
	if p.match([]TokenType{LeftParen}) {
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
			fallthrough
		case Fun:
			fallthrough
		case Var:
			fallthrough
		case For:
			fallthrough
		case If:
			fallthrough
		case While:
			fallthrough
		case Print:
			fallthrough
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
