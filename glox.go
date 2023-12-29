package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	// https://man.freebsd.org/cgi/man.cgi?query=sysexits
	// The command was used	incorrectly.
	SysexitsUsage = 64
	// The input data was incorrect	in some	way.
	SysexitsDataError = 65
	// An internal software	error has been detected.
	SysexitsUsageSoftware = 70
)

var intepreter = NewInterpreter()

// Whether we've encountered an error.
var hadError bool

// Whether we've encountered a runtime error.
var hadRuntimeError bool

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %q: %v\n", path, err)
		os.Exit(SysexitsUsageSoftware)
	}
	run(string(bytes))

	if hadError {
		println("Encountered error. Exiting.")
		// Indicate an error in the exit code.
		os.Exit(SysexitsDataError)
	}
	if hadRuntimeError {
		println("Encountered runtime error. Exiting.")
		os.Exit(SysexitsUsageSoftware)
	}
}

func runPrompt() {
	// Seems docs recommend using Scanner vs. Reader for reading line-by-line.
	reader := bufio.NewReader(os.Stdin)

	for {
		print("> ")
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			os.Exit(1)
		}
		if bytes == nil {
			break
		}
		run(string(bytes))
		// Don't kill the session if the user makes an error.
		hadError = false
	}
}

func run(source string) {
	s := NewScanner(source)
	tokens := s.ScanTokens()
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	// TODO: Find the best way to handle errors discovered.
	if err != nil {
		hadError = true
	}
	if hadError {
		return
	}
	intepreter.Interpret(statements)
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

func PrintDetailedError(token Token, message string) {
	if token.TokenType == EOF {
		report(token.Line, "at end", message)
	} else {
		report(token.Line, "at '"+token.Lexeme+"'", message)
	}
}

func PrintRuntimeError(err RuntimeError) {
	println(fmt.Sprintf("%s\n[line %d]", err.Message, err.Token.Line))
	hadRuntimeError = true
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
