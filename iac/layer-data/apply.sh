#!/bin/bash

REGION="europe-west1"
GCP_PROJECT="slavayssiere-sandbox"

terraform apply \
    --var "region=$REGION" \
    -auto-approve

BT_INSTANCE="test-instance"
BT_TABLE="test-table"

gcloud beta bigtable instances create $BT_INSTANCE \
   --cluster=$BT_INSTANCE \
   --cluster-zone="europe-west1-b" \
   --display-name=$BT_INSTANCE \
   --cluster-num-nodes=3

# Sample code to bootstrap BigTable structure and inserting a sample value
cbt -instance $BT_INSTANCE createtable $BT_TABLE
cbt -instance $BT_INSTANCE createfamily $BT_TABLE ms

gcloud beta scheduler jobs create pubsub aggregator-stats-call \
    --description="Launch job to create aggregas" \
    --schedule="*/5 * * * *" \
    --topic="projects/slavayssiere-sandbox/topics/aggregator-queue" \
    --message-body-from-file=./message-stats-aggregator.json


gcloud beta scheduler jobs create pubsub aggregator-dataset-call \
    --description="Launch job to create dataset" \
    --schedule="0 0 * * *" \
    --topic="projects/slavayssiere-sandbox/topics/aggregator-queue" \
    --message-body-from-file=./message-stats-dataset.json

