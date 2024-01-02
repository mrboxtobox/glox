package main

import "fmt"

// AST representation of a token.
type Token struct {
	TokenType TokenType
	// Raw substring that forms the token.
	Lexeme string
	// The parsed representation of literal values e.g. float64, string.
	Literal interface{}
	// Line where the token was scanned.
	Line int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{type=%s, lexeme=%s, literal=%v}", t.TokenType, t.Lexeme, t.Literal)
}
