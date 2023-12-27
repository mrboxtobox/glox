BINARY_NAME=glox

all:
	# Cleaning past artifacts...
	go clean
	rm -f ./glox
	rm -f ./generate_ast

	# Build the AST generator...
	# go build -o ./generate_ast tool/generate_ast.go

	# Generate expr.go...
	# ./generate_ast .

	# Build interpreter...
	go build -o ./glox *.go

	# Run interpreter...
	./glox

build:
	go build -o ./glox *.go

run:
	./glox

build_and_run: build run

genast:
	# rm generate_ast
	go build -o ./generate_ast tool/generate_ast.go
	./generate_ast .

clean:
	go clean
	rm ./glox
	rm ./generate_ast