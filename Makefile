BINARY_NAME=glox

all:
	# Removing existing artifacts...
	go clean
	rm -rf ./bin

	# Building the AST generator...
	go build -o ./bin/genast tool/*.go

	# Generating AST files...
	./bin/genast .

	# Building the interpreter...
	go build -o ./bin/glox *.go

run: all
	# Running the interpreter...
	./bin/glox

example: all
	##############################
	# Running the interpreter... #
	##############################
	./bin/glox examples/fibonacci.g
