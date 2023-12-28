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

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err == nil {
		return expr
	}
	// Use type assertion to determine error type.
	// See https://go.dev/tour/methods/15.
	if _, ok := err.(ParseError); ok {
		return nil
	}
	// TODO: We should never end in this state so, panic.
	panic(err)
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
		expr = Binary{expr, operator, right}
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
		expr = Binary{expr, operator, right}
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
		expr = Binary{expr, operator, right}
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
		expr = Binary{expr, operator, right}
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
		return Unary{operator, right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match([]TokenType{False}) {
		return Literal{false}, nil
	}
	if p.match([]TokenType{True}) {
		return Literal{true}, nil
	}
	if p.match([]TokenType{Nil}) {
		return Literal{nil}, nil
	}
	if p.match([]TokenType{Number, String}) {
		return Literal{p.previous().Literal}, nil
	}
	if p.match([]TokenType{LeftParen}) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RightParen, "Expect ')' after expression but got."); err != nil {
			return nil, err
		}
		return Grouping{expr}, nil
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
