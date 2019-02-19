#!/bin/bash

REGION="europe-west1"
GCP_PROJECT="slavayssiere-sandbox"

gcloud -q beta compute routers nats delete nat-$REGION --router=nat-$REGION
gcloud -q beta compute routers delete nat-$REGION

touch empty-file
gcloud dns record-sets import -z private-dns-zone --delete-all-existing empty-file
rm empty-file
gcloud -q beta dns managed-zones delete private-dns-zone

terraform destroy \
    -auto-approve

## for injectors
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.publisher >> /dev/null
## for consumer
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
## for normalizers
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.publisher >> /dev/null
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/automl.admin >> /dev/null
## for aggregators
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner >> /dev/null
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
## for sa-pubsub-bigtable
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com --role roles/bigtable.admin >> /dev/null
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com --role roles/bigtable.user >> /dev/null
## for aggregator
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-datastore@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-datastore@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner >> /dev/null
##for bastion
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-bastion@$GCP_PROJECT.iam.gserviceaccount.com --role roles/container.clusterAdmin >> /dev/null
##for sse
gcloud  -q projects remove-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-sse@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.editor >> /dev/null

gcloud -q iam service-accounts delete "sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-pubsub-datastore@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-bastion@$GCP_PROJECT.iam.gserviceaccount.com"
gcloud -q iam service-accounts delete "sa-sse@$GCP_PROJECT.iam.gserviceaccount.com"
