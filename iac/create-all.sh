#!/bin/bash

source ../env.sh

cd layer-base
./apply.sh
cd -

cd layer-bastion
./apply.sh
cd -

cd layer-kubernetes
./apply.sh
cd -

cd layer-data
./apply.sh
cd -

cd layer-services
./apply.sh
cd -

cd ../visualizer
./apply.sh
cd -

cd layer-istio-lb-http
./apply.sh
cd -

ASSET_DOMAIN="assets.gcp-wescale.slavayssiere.fr"

gsutil mb gs://$ASSET_DOMAIN

gsutil cp ../app-sse/src/templates/index.html gs://$ASSET_DOMAIN
gsutil cp ../app-sse/src/templates/twitter.png gs://$ASSET_DOMAIN
gsutil iam ch allUsers:objectViewer gs://$ASSET_DOMAIN
gsutil web set -m index.html -e 404.html gs://$ASSET_DOMAIN

cd ../deployment
./deploy.sh
cd -

cd ../functions
./apply.sh
cd -
