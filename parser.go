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
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match([]TokenType{Print}) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(Semicolon, "Expect ';' after value.")
	return PrintStmt{value}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(Semicolon, "Expect ';' after expression.")
	return ExpressionStmt{expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
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
