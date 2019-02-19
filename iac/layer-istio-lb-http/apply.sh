#!/bin/bash

source ../../env.sh

terraform apply -auto-approve

kubectl apply -f istio-ingress.yaml

until (gcloud compute health-checks list | grep k8s)
do
    echo "Wait for healthcheck..."
    sleep 10
done

hcs=$(gcloud compute health-checks list --protocol=HTTP | grep k8s | cut -d ' ' -f1)
porthc=$(kubectl -n istio-system get svc admin-ingressgateway -o jsonpath='{.spec.healthCheckNodePort}')

while read -r line; do   
    echo "change health for '$line' on port: $porthc "
    gcloud compute health-checks update http $line --request-path=/healthz --port=$porthc
done <<< "$hcs"
