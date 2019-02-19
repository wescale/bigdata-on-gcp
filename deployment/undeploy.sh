#!/bin/bash

GCP_PROJECT="slavayssiere-sandbox"

kubectl delete -f .

################################ Injectors ################################
source ../env.sh
kubectl delete secret twitter-secrets \
    -n injectors

kubectl delete secret sa-pubsub-publisher \
    -n injectors

################################ Normalizers ################################
kubectl delete secret sa-pubsub-full \
    -n normalizers

################################ Consumers ################################
kubectl delete secret sa-pubsub-subscriber \
    -n consumers

kubectl delete secret sa-pubsub-bigtable \
    -n consumers


################################ Aggregators ################################
kubectl delete secret sa-aggregator \
    -n aggregators
