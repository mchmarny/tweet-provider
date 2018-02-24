# twitter-to-pubsub-event-pump [![Build Status](https://travis-ci.org/mchmarny/twitter-to-pubsub-event-pump.svg?branch=master)](https://travis-ci.org/mchmarny/twitter-to-pubsub-event-pump)

Utility to subscribe to Twitter events with query and pump events into Google PubSub topic.

## Setup

### Twitter 

To run this app you will need to obtain your personal Twitter app Consume and OAuth secrets (`Consumer Key`,
`Consumer Secret`,`OAuth Access Token`,`OAuth Access Token Secret`) Good instructions on how to obtain these are located [here](https://iag.me/socialmedia/how-to-create-a-twitter-app-in-8-easy-steps/)

Once you obtain these four secrets from Twitter, you will need to pass these as arguments on each execution or you can define them as environment variables:

```shell
export T_CONSUMER_KEY="***"
export T_CONSUMER_SECRET="***"
export T_ACCESS_TOKEN="***"
export T_ACCESS_SECRET="***"
```

### GCP

If you don't already have GCP account, you can run this entire app using the Google Cloud Platform (GCP) [free tier](https://cloud.google.com/free/). Once you create project, you will need to pass the `GCLOUD_PROJECT` argument on each execution or you can define it as environment variables like this:

```shell 
export GCLOUD_PROJECT="YOUR_PROJECT_NAME"
```

### GCP CLI

If you don't already have `gcloud`, you can find instructions on how to download and install the GCP SDK [here](https://cloud.google.com/sdk/)


#### Service Account 

You will need to set up GCP authentication using service account. You can find instructions how to do this [here](https://cloud.google.com/video-intelligence/docs/common/auth#set_up_a_service_account). After you download your service account file you will need to define

```shell 
export GOOGLE_APPLICATION_CREDENTIALS=<path_to_service_account_file>
```

#### Create PubSub topic

You can create GCP PubSub topic `gcloud` by executing the following command:

```shell 
gcloud pubsub topics create tweets
```

## Build

To first build the app you can execute first `make dep` which will assure your environment has the necessary dependencies. Alternatively you can restore the app dependencies 

```shell
go get github.com/tools/godep
godep restore
```

and then `make build` to build the app or alternatively you can run the build command directly

```shell
go build -v -o bin/tpump
```

## Run 

You can run the app (assuming it's already built) using the following command. You can use the `serverless` sample or use your own query. The twitter search query operators are outlined [here](https://developer.twitter.com/en/docs/tweets/search/guides/standard-operators)

```shell
bin/tpump --query="serverless OR faas OR openwhisk OR openfaas OR lambda"
```

This command will search Twitter for tweets matching your query and push them one by one into GCP PubSub topic. 
