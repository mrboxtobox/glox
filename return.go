package main

import "fmt"

type FunctionReturn struct {
	Value any
}

func (fr FunctionReturn) Error() string {
	return fmt.Sprintf("Error: %v", fr.Value)
}
