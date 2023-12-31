package main

type RuntimeError struct {
	Token   Token
	Message string
}

func (e RuntimeError) Error() string {
	return "RuntimeError: " + e.Message
}

func LogAndReturnError(token Token, message string) RuntimeError {
	PrintDetailedError(token, message)
	return RuntimeError{token, message}
}
