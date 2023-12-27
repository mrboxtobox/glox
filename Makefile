BINARY_NAME=glox

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