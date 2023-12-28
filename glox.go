package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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

var intepreter = Interpreter{}

// Whether we've encountered an error.
var hadError bool

// Whether we've encountered a runtime error.
var hadRuntimeError bool

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		os.Exit(SysexitsUsageSoftware)
	}
	run(string(bytes))

	if hadError {
		// Indicate an error in the exit code.
		os.Exit(SysexitsDataError)
	}
	if hadRuntimeError {
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
	expr := parser.Parse()
	if hadError {
		return
	}
	printer := AstPrinter{}
	printed, err := printer.print(expr)
	if err != nil {
		log.Fatal("AstPrinter failed to print: ", err)
	}
	fmt.Printf("%v\n", printed)
	intepreter.Interpret(expr)
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
	args := flag.Args()
	if len(args) > 1 {
		println("Usage: glox [script]")
		os.Exit(SysexitsUsage)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
