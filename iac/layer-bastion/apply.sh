#!/bin/bash

REGION="europe-west1"
GCP_PROJECT="slavayssiere-sandbox"

terraform apply \
    --var "region=$REGION" \
    -auto-approve

    