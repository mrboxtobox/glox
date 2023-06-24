BINARY_NAME=glox

build:
	go build -o ./${BINARY_NAME} glox.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ./${BINARY_NAME}