#!/bin/bash

gcloud -q beta scheduler jobs delete aggregator-stats-call
gcloud -q beta scheduler jobs delete aggregator-dataset-call

cbt deleteinstance "test-instance"

gcloud dataflow jobs run delete-datastore \
    --gcs-location gs://dataflow-templates/latest/Datastore_to_Datastore_Delete \
    --parameters \
datastoreReadGqlQuery="SELECT * FROM aggregas",\
datastoreReadProjectId="slavayssiere-sandbox",\
datastoreDeleteProjectId="slavayssiere-sandbox"

terraform destroy \
    -auto-approve
