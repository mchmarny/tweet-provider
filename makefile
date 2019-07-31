.PHONY: app client service

all: test

# DEV
run:
	go run *.go -v

query:
	curl -H "Content-Type: application/json" \
		 -X POST -d '{"query":"serverless"}' http://localhost:8080/

# BUILD
mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
      --project cloudylabs-public \
	  --tag gcr.io/cloudylabs-public/twitter-to-pubsub-event-pump:0.4.1

# DEPLOY
auth: meta
	gcloud projects add-iam-policy-binding cloudylabs \
    	--member="serviceAccount:service-${PROJECT_NUM}@gcp-sa-pubsub.iam.gserviceaccount.com" \
    	--role=roles/iam.serviceAccountTokenCreator

meta:
	PROJECT=$(gcloud config get-value project)
	PROJECT_NUM=$(gcloud projects list --filter="${PROJECT}" --format="value(PROJECT_NUMBER)")

service:
	gcloud beta run deploy twitter-query \
		--image=gcr.io/cloudylabs-public/twitter-to-pubsub-event-pump:0.4.1 \
		--region=us-central1 \
		--timeout=15m

serviceless:
	gcloud beta run services delete twitter-query

sa:
	gcloud iam service-accounts create cr-demo-sa \
    	--display-name "Cloud Run Service Invoker"

	gcloud beta run services add-iam-policy-binding cr-demo \
		--member=serviceAccount:cr-demo-sa@cloudylabs.iam.gserviceaccount.com \
		--role=roles/run.invoker