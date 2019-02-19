#!/bin/bash

test_tiller_present() {
    kubectl get pod -n kube-system -l app=helm,name=tiller | grep Running | wc -l | tr -d ' '
}

apply_kubectl() {
    cd $1
    kubectl apply -f .
    cd ..
}

gcloud container clusters get-credentials test-cluster --region europe-west1
username=$(gcloud config get-value account)
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$username
    
# install prometheus operator
kubectl apply -f helm/rbac.yaml
helm init  --service-account tiller

test_tiller=$(test_tiller_present)
while [ $test_tiller -lt 1 ]; do
    echo "Wait for Tiller: $test_tiller"
    test_tiller=$(test_tiller_present)
    sleep 1
done

sleep 10

kubectl create ns monitoring
helm install --name prometheuses stable/prometheus-operator --namespace monitoring

apply_kubectl "external-dns"

if [ ! -d "istio-1.0.5" ]; then
    wget https://github.com/istio/istio/releases/download/1.0.5/istio-1.0.5-osx.tar.gz
    tar -xvf istio-1.0.5-osx.tar.gz
    rm istio-1.0.5-osx.tar.gz
fi

cd istio-1.0.5
    helm install install/kubernetes/helm/istio --name istio --namespace istio-system -f ../istio/values-istio-1.0.5.yaml
cd -

kubectl -n istio-system annotate svc iap-ingressgateway beta.cloud.google.com/backend-config="{\"ports\": {\"http2\" :\"config-iap\"}, \"default\": \"config-iap\"}" --overwrite=true
kubectl -n istio-system annotate svc public-ingressgateway beta.cloud.google.com/backend-config="{\"ports\": {\"http2\" :\"config-default\"}, \"default\": \"config-default\"}" --overwrite=true

# configure cloud IAP
source ../../env.sh
kubectl -n istio-system create secret generic my-oauth-secret \
	--from-literal=client_id=$CLIENT_ID \
    --from-literal=client_secret=$CLIENT_SECRET
kubectl apply -f cloud-service/backend-config.yaml

apply_kubectl "monitoring"
