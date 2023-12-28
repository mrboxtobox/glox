BINARY_NAME=glox

all:
	# Cleaning past artifacts...
	go clean
	rm -f ./glox
	rm -f ./generate_ast

	# Build the AST generator...
	go build -o ./genast tool/*.go

	# Generate expr.go...
	./genast .

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
	go build -o ./genast tool/generate_ast.go
	./genast .

clean:
	go clean
	rm -f ./glox
	rm -f ./genast