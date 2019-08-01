# tweet-provider

Simple Twitter search service with PubSub result publishing and Firestore managed state

## Pre-requirements

### Twitter

To run this app you will need to obtain your personal Twitter app Consumer and OAuth secrets (`Consumer Key`,
`Consumer Secret`,`OAuth Access Token`,`OAuth Access Token Secret`) Good instructions on how to obtain these are located [here](https://iag.me/socialmedia/how-to-create-a-twitter-app-in-8-easy-steps/)

Once you obtain these four secrets from Twitter, you will need to pass these as arguments on each execution or you can define them as environment variables:

```shell
export T_CONSUMER_KEY="***"
export T_CONSUMER_SECRET="***"
export T_ACCESS_TOKEN="***"
export T_ACCESS_SECRET="***"
```

### GCP CLI

If you don't already have `gcloud`, you can find instructions on how to download and install the GCP SDK [here](https://cloud.google.com/sdk/)

## Setup

### Build Container Image

Cloud Run uses container images so let's start bu building one....

> PS. you can review each one of the provided scripts for complete gcloud command

```shell
bin/image
```

### Service Account and IAM Policies

Next we will need to create a service account and assign it all the necessary IAM roles...

```shell
bin/account
```

### Cloud Run Service

Once you have the service account we can now deploy a new service and set it to run under that account.

> Note, the deployed service will require authentication and run under the the service account identity we configured in the previous step

```shell
bin/service
```

### Cloud Schedule

The Cloud Run service will search Twitter for provided query so now we just have to create a Cloud Schedule to execute the service every 10 min.

> Note, the Cloud Run service stores the state for each query to that can start searches from the maximum tweet returned from the previous query.

```shell
bin/schedule
```

## Cleanup

To cleanup all resources created by this sample execute

```shell
bin/cleanup
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.