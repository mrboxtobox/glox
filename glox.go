package main

import (
	"bufio"
	"flag"
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

// Whether we've encountered an error.
var hadError bool

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
	s := newScanner(source)
	tokens := s.scanTokens()

	for _, token := range tokens {
		println(token.String())
	}
}

// Minimal error reporting.
func printErr(line int, message string) {
	report(line, "", message)
}

// TODO: Extend this to show users the offending line and point to the column.
func report(line int, where, message string) {
	println(fmt.Sprintf("[%d] Error %s: %s", line, where, message))
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
