#!/bin/bash

GCP_PROJECT="slavayssiere-sandbox"

gcloud -q functions delete laststat
gcloud -q functions delete getstat

gcloud -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner

gcloud -q iam service-accounts delete "sa-cloudfunction@$GCP_PROJECT.iam.gserviceaccount.com"

