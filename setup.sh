#!/bin/bash



# populate twitter secrets
kubectl create secret generic tpump-tw-key --from-literal=T_CONSUMER_KEY=$T_CONSUMER_KEY
kubectl create secret generic tpump-tw-secret --from-literal=T_CONSUMER_SECRET=$T_CONSUMER_SECRET
kubectl create secret generic tpump-tw-token --from-literal=T_ACCESS_TOKEN=$T_ACCESS_TOKEN
kubectl create secret generic tpump-tw-access --from-literal=T_ACCESS_SECRET=$T_ACCESS_SECRET

# populate app secrets
kubectl create secret generic tpump-gcloud-project --from-literal=GCLOUD_PROJECT=$GCLOUD_PROJECT
kubectl create secret generic s0-demo-sa --from-file ~/.gcp-keys/s9-demo-sa.json

# setup query
kubectl create configmap tpump-query --from-literal=T_QUERY="serverless OR faas OR openwhisk OR openfaas OR lambda"

# deploy
kubectl create -f tpump.yaml