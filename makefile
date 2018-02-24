
# Go parameters
BINARY_NAME=tpump

all: test build

build:
	go build -v -o bin/$(BINARY_NAME) 

test:
	go test -v ./...

clean:
	go clean
	rm -f bin/$(BINARY_NAME)

run: build
	bin/$(BINARY_NAME) --query="serverless OR faas OR openwhisk OR openfaas OR lambda"

deps:
	go get github.com/tools/godep
	godep restore

