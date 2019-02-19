#!/bin/bash

cd layer-base
terraform init \
    -backend-config="region=europe-west1"
cd -

cd layer-bastion
terraform init \
    -backend-config="region=europe-west1"
cd -

cd layer-kubernetes
terraform init \
    -backend-config="region=europe-west1"
cd -