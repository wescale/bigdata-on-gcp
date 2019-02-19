#!/bin/bash

REGION="europe-west1"
SUBNET_EUROPE1_CIDR="192.168.0.0/20"
GCP_PROJECT="slavayssiere-sandbox"

gcloud config set project $GCP_PROJECT
gcloud config set compute/region $REGION

terraform apply \
    --var "region=europe-west1" \
    -auto-approve

gcloud -q beta compute routers create nat-$REGION  \
    --network demo-net \
    --region $REGION

gcloud -q beta compute routers nats create nat-$REGION \
    --router-region $REGION \
    --router nat-$REGION \
    --nat-all-subnet-ip-ranges \
    --auto-allocate-nat-external-ips

gcloud -q beta dns managed-zones create private-dns-zone \
    --dns-name="gcp.wescale" \
    --description="A private zone" \
    --visibility=private \
    --networks=demo-net


gcloud iam service-accounts create "sa-pubsub-publisher" --display-name "SA for pubsub publish apps"
gcloud iam service-accounts create "sa-pubsub-subscriber" --display-name "SA for pubsub publish apps"
gcloud iam service-accounts create "sa-pubsub-full" --display-name "SA for pubsub publish apps"
gcloud iam service-accounts create "sa-aggregator" --display-name "SA for aggregator apps"
gcloud iam service-accounts create "sa-pubsub-bigtable" --display-name "SA for pubsub and bigtable apps"
gcloud iam service-accounts create "sa-pubsub-datastore" --display-name "SA for pubsub and datastore apps"
gcloud iam service-accounts create "sa-bastion" --display-name "SA for bastion"
gcloud iam service-accounts create "sa-sse" --display-name "SA for sse"



## for injectors
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.publisher >> /dev/null
## for consumer
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
## for normalizers
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.publisher >> /dev/null
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com --role roles/automl.admin >> /dev/null
## for aggregators
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner >> /dev/null
## for sa-pubsub-bigtable
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com --role roles/bigtable.admin >> /dev/null
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com --role roles/bigtable.user >> /dev/null
## for aggregator
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-datastore@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.subscriber >> /dev/null
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-pubsub-datastore@$GCP_PROJECT.iam.gserviceaccount.com --role roles/datastore.owner >> /dev/null
##for bastion
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-bastion@$GCP_PROJECT.iam.gserviceaccount.com --role roles/container.clusterAdmin >> /dev/null
##for sse
gcloud -q projects add-iam-policy-binding $GCP_PROJECT --member serviceAccount:sa-sse@$GCP_PROJECT.iam.gserviceaccount.com --role roles/pubsub.editor >> /dev/null

gcloud iam service-accounts keys create ../sa-pubsub-publisher.json \
  --iam-account sa-pubsub-publisher@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-pubsub-subscriber.json \
  --iam-account sa-pubsub-subscriber@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-pubsub-full.json \
  --iam-account sa-pubsub-full@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-aggregator.json \
  --iam-account sa-aggregator@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-pubsub-bigtable.json \
  --iam-account sa-pubsub-bigtable@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-pubsub-datastore.json \
  --iam-account sa-pubsub-datastore@$GCP_PROJECT.iam.gserviceaccount.com

gcloud iam service-accounts keys create ../sa-sse.json \
  --iam-account sa-sse@$GCP_PROJECT.iam.gserviceaccount.com

