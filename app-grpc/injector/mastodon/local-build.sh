#!/bin/bash

version=$1

export GCP_PROJECT=$(gcloud config get-value project)
export TOPIC_NAME="projects/$GCP_PROJECT/topics/mastodon-raw"
export SECRET_PATH="/Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-publisher.json"

docker build -t eu.gcr.io/$GCP_PROJECT/mastodon-injector:$version .
docker push eu.gcr.io/$GCP_PROJECT/mastodon-injector:$version
