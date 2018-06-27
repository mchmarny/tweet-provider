
# Go parameters
BINARY_NAME=tpump
GCP_PROJECT_NAME=knative-samples

all: test
build:
	go build -o ./bin/$(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)

run: build
	bin/$(BINARY_NAME) --query="serverless OR faas OR openwhisk OR openfaas OR lambda"

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure

gcr:
	gcloud container builds submit --project=$(GCP_PROJECT_NAME) --tag gcr.io/$(GCP_PROJECT_NAME)/$(BINARY_NAME):latest .
