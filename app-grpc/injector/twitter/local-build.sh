# Localy do:

version=$1

export GCP_PROJECT=$(gcloud config get-value project)
export TOPIC_NAME="projects/$GCP_PROJECT/topics/twitter-raw"
export SECRET_PATH="/Users/slavayssiere/Code/slavayssiere-sandbox-gcp/iac/sa-pubsub-publisher.json"

docker build -t eu.gcr.io/$GCP_PROJECT/twitter-injector:$version .
docker push eu.gcr.io/$GCP_PROJECT/twitter-injector:$version
