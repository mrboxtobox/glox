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

run: all
	# Run interpreter...
	./glox
