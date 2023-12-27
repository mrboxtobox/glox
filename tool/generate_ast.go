package tool

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// defineType generates the AST for a single type.
func defineType(file *os.File, baseName string, structName string, fields string) {
	file.WriteString("type " + structName + " struct {\n")
	// Fields.
	for _, field := range strings.Split(fields, ",") {
		file.WriteString(" " + strings.TrimSpace(field))
		file.WriteString("\n")
	}
	file.WriteString("}\n\n")
	file.WriteString("\n  func " + "(expr " + structName + ")" + "accept(visitor Visitor) any{\n")
	file.WriteString("  return visitor.visit" + structName + baseName + "(expr)\n")
	file.WriteString("}\n")
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

	file.WriteString("package main\n\n")
	// The Visitor.
	defineVisitor(file, baseName, types)

	file.WriteString("type " + baseName + " interface {\n")
	file.WriteString("  accept(visitor Visitor) any")
	file.WriteString("}\n\n")

	// Add the AST Classes.
	for _, t := range types {
		parts := strings.Split(t, ":")
		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(file, baseName, structName, fields)
	}
	// The accept() method.
	// file.WriteString("\n  func " + "(" + baseName + ")" + "accept(visitor Visitor) any")
}

// TODO: Define this.
func defineVisitor(file *os.File, baseName string, types []string) {
	file.WriteString("type Visitor interface {\n")
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		funcName := "  visit" + typeName + baseName
		param := strings.ToLower(baseName) + " " + typeName
		file.WriteString(funcName + "(" + param + ") any\n")
	}
	file.WriteString("}\n\n")
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	// Define the AST.
	dir := os.Args[1]
	defineAst(dir, "Expr", []string{
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator Token, right Expr",
	})
	formatFiles()
}
