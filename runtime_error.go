package main

type RuntimeError struct {
	Token   Token
	Message string
}

func (e RuntimeError) Error() string {
	return "RuntimeError: " + e.Message
}
