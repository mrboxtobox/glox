package main

import "fmt"

interface expr {}

type Binary    struct {
	left Expr
	operator tokenType
	right Expr
}
type Grouping  struct {
	expression Expr
}
type Literal   struct {
	value interface{}
}
type Unary     struct {
	operator tokenType
	right Expr
}
}
