package main

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{type=%s, lexeme=%s, literal=%v}", t.TokenType, t.Lexeme, t.Literal)
}
