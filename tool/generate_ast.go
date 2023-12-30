package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO: Check error status in WriteString.
// defineType generates the AST for a single type.
// Wrap prefix with basename since we can't have nested classes?
func defineType(file *os.File, baseName string, structName string, fields string) {
	WriteStringOrDie(file, "type "+structName+baseName+" struct {\n")
	// Fields.
	for _, field := range strings.Split(fields, ",") {
		WriteStringOrDie(file, " "+strings.TrimSpace(field))
		WriteStringOrDie(file, "\n")
	}
	WriteStringOrDie(file, "}\n\n")
	WriteStringOrDie(file, "\n  func "+"(expr "+structName+baseName+") Accept"+baseName+"(visitor "+baseName+"Visitor) (any, error) {\n")
	WriteStringOrDie(file, "  return visitor.Visit"+structName+baseName+"(expr)\n")
	WriteStringOrDie(file, "}\n")
}

// defineAst generates the AST for a set of supported types.
func defineAst(dir string, baseName string, types []string) {
	path := dir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Unable to create file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	WriteStringOrDie(file, "package main\n\n")
	// The Visitor.
	defineVisitor(file, baseName, types)

	WriteStringOrDie(file, "type "+baseName+" interface {\n")
	WriteStringOrDie(file, "  Accept"+baseName+"(visitor "+baseName+"Visitor) (any, error)")
	WriteStringOrDie(file, "}\n\n")

	// Add the AST Classes.
	for _, t := range types {
		parts := strings.Split(t, ":")
		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(file, baseName, structName, fields)
	}
}

// We create separate visitors to avoid conflicts.
func defineVisitor(file *os.File, baseName string, types []string) {
	WriteStringOrDie(file, "type "+baseName+"Visitor interface {\n")
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		funcName := "  Visit" + typeName + baseName
		param := strings.ToLower(baseName) + " " + typeName + baseName
		WriteStringOrDie(file, funcName+"("+param+") (any, error)\n")
	}
	WriteStringOrDie(file, "}\n\n")
}

func formatFiles() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path: %v\n", err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fmt.Printf("Formatting: %s\n", path)
			cmd := exec.Command("gofmt", "-w", path)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error running gofmt: %v\n", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
}

func WriteStringOrDie(file *os.File, str string) {
	if _, err := file.WriteString(str); err != nil {
		panic((fmt.Sprintf("Error while writing to file: %v\n", err)))
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	// Define the AST.
	dir := os.Args[1]
	defineAst(dir, "Expr", []string{
		"Assign   : Name Token, Value Expr",
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
		"Logical  : Left Expr, Operator Token, Right Expr",
		"Unary    : Operator Token, Right Expr",
		"Variable : Name Token",
	})

	defineAst(dir, "Stmt", []string{
		"Block      : Statements []Stmt",
		"Expression : Expression Expr",
		"If         : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print      : Expression Expr",
		"Var        : Name Token, Initializer Expr",
		"While      : Condition Expr, Body Stmt",
	})
	formatFiles()
}
