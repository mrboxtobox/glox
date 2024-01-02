package main

import (
	"bufio"
	"fmt"
	"os"
)

// See https://man.freebsd.org/cgi/man.cgi?query=sysexits.
const (
	// The command was used incorrectly.
	SysexitsUsage = 64
	// The input data was incorrect in some way.
	SysexitsDataError = 65
	// An internal software error has been detected.
	SysexitsUsageSoftware = 70
)

var intepreter = NewInterpreter()

// Whether we've encountered an error.
var hadError bool

// Whether we've encountered a runtime error.
var hadRuntimeError bool

func PrintDetailedError(token Token, message string) {
	if token.TokenType == EOFToken {
		report(token.Line, "at end", message)
	} else {
		report(token.Line, "at '"+token.Lexeme+"'", message)
	}
}

func PrintRuntimeError(err RuntimeError) {
	println(fmt.Sprintf("[line %d]%s\n", err.Token.Line, err.Message))
	hadRuntimeError = true
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %q: %v\n", path, err)
		os.Exit(SysexitsUsageSoftware)
	}
	err = run(string(bytes))
	switch typedErr := err.(type) {
	case RuntimeError:
		fmt.Printf("[line %d]%s\n", typedErr.Token.Line, typedErr.Message)
		println("Encountered a runtime error. Exiting.")
		os.Exit(SysexitsUsageSoftware)
	default:
		fmt.Printf("Encountered an unexpected error (%T). Exiting.", typedErr)
		os.Exit(SysexitsUsageSoftware)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		print("> ")
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Encountered error wile reading input: %v.\n", err)
			os.Exit(1)
		}
		if bytes == nil {
			break
		}
		// Don't kill the session if the user makes an error.
		_ = run(string(bytes))
	}
}

func run(source string) error {
	s := NewScanner(source)
	tokens := s.ScanTokens()
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}
	resolver := NewResolver(intepreter)
	// Stop if there was a resolution error.
	if _, err := resolver.resolveAll(statements); err != nil {
		return err
	}
	if hadError {
		return fmt.Errorf("Encountered an error")
	}
	if err := intepreter.Interpret(statements); err != nil {
		return err
	}
	return nil
}

// Minimal error reporting.
func printErr(line int, message string) {
	report(line, "", message)
}

// TODO: Extend this to show users the offending line and point to the column.
func report(line int, where, message string) {
	println(fmt.Sprintf("[%d] Error %s: %s", line, where, message))
	hadError = true
}

func main() {
	args := os.Args
	if len(args) > 2 {
		println("Usage: glox [script]")
		os.Exit(SysexitsUsage)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}
