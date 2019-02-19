#!/bin/bash

REGION="europe-west1"
MYIP=$(curl ifconfig.me)

terraform apply \
    --var "region=$REGION" \
    --var "myip=$MYIP" \
    -auto-approve
