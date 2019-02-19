#!/bin/bash

sudo apt-get install kubectl


REGION="europe-west1"
NAME_CLUSTER="test-cluster"

echo "gcloud config set compute/region $REGION" | sudo tee --append /home/admin/.bashrc  > /dev/null
echo "gcloud container clusters get-credentials $NAME_CLUSTER --region $REGION" | sudo tee --append /home/admin/.bashrc  > /dev/null
echo "alias ll='ls -lisa'" | sudo tee --append /home/admin/.bashrc  > /dev/null

