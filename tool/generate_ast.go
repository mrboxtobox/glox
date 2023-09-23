package main

import (
	"fmt"
	"os"
	"strings"
)

// defineType generates the AST for a single type.
func defineType(file *os.File, baseName string, structName string, fields string) {
	file.WriteString("type " + structName + " struct {\n")
	// Fields.
	for _, field := range strings.Split(fields, ",") {
		file.WriteString(strings.TrimSpace(field))
	}
	file.WriteString("}\n")
}

// defineAst generates the AST for a set of supported types.
func defineAst(dir string, baseName string, types []string) {
	path := dir + "/" + baseName + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Unable to create file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString("package main\n\n")
	file.WriteString("import \"fmt\"\n\n")
	file.WriteString("interface " + baseName + " {}\n\n")
	// Add the AST Classes.
	for _, t := range types {
		parts := strings.Split(t, ":")
		structName := parts[0]
		fields := parts[1]
		defineType(file, baseName, structName, fields)
	}

	file.WriteString("}\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	// Define the AST.
	dir := os.Args[1]
	defineAst(dir, "expr", []string{
		"Binary   : left Expr, operator tokenType, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator tokenType, right Expr",
	})
}
