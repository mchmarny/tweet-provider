.PHONY: app client service

all: test

# DEV
run:
	go run *.go -v

# BUILD
mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
      --project cloudylabs-public \
	  --tag gcr.io/cloudylabs-public/twitter-to-pubsub-event-pump:0.4.1

query:
	curl -H "Content-Type: application/json" \
		 -X POST -d '{"query":"serverless"}' http://localhost:8080/