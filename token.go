package main

import "fmt"

type token struct {
	tokenType tokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t token) String() string {
	return fmt.Sprintf("Token{type=%s, lexeme=%s, literal=%v}", t.tokenType, t.lexeme, t.literal)
}
