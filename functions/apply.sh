#!/bin/bash

GCP_PROJECT="slavayssiere-sandbox"

gcloud iam service-accounts create "sa-cloudfunction" --display-name "SA for bastion"
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner >> /dev/null

cd laststat
gcloud alpha functions deploy laststat \
    --entry-point LastStat \
    --runtime go111 \
    --trigger-http \
    --service-account "sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com"
cd -

cd getstat
gcloud alpha functions deploy getstat \
    --entry-point GetStat \
    --runtime go111 \
    --trigger-http \
    --service-account "sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com"
cd -

